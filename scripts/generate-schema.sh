#!/bin/bash
# scripts/extract-schema.sh

# Default values
MIGRATIONS_DIR="../internal/database/migrations"
DEFAULT_FILE="schemas.sql"
OUTPUT_FILE="$MIGRATIONS_DIR/$DEFAULT_FILE"

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            # Check if user provided full path or just filename
            if [[ "$2" == */* ]]; then
                # User provided full path
                OUTPUT_FILE="$2"
            else
                # User provided just filename, prepend migrations dir
                OUTPUT_FILE="$MIGRATIONS_DIR/$2"
            fi
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [-o output_file.sql]"
            echo "Options:"
            echo "  -o, --output    Specify output file name or path"
            echo "                  (default: $MIGRATIONS_DIR/schemas.sql)"
            echo "  -h, --help      Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h for help"
            exit 1
            ;;
    esac
done

# Create migrations directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Clear the output file
> "$OUTPUT_FILE"

echo "🔍 Searching for schemas..."
echo "📁 Output file: $OUTPUT_FILE"

find ../internal/api/domain -name "model.go" | while read FILE; do
    echo "📄 Processing: $FILE"
    
    # Read the entire file
    CONTENT=$(cat "$FILE")
    
    # Find everything between backticks after "var.*Schema.*="
    SQL=$(echo "$CONTENT" | perl -ne '
        if (/var.*[Ss]chema.*=\s*`(.*)/) {
            $in_sql = 1;
            $sql = $1;
            if ($sql =~ /`/) {
                $sql =~ s/`.*//;
                print $sql;
                $in_sql = 0;
            } else {
                print $sql . "\n";
            }
        } elsif ($in_sql) {
            if (/`/) {
                s/`.*//;
                print $_;
                $in_sql = 0;
            } else {
                print $_;
            }
        }
    ')
    
    if [ -n "$SQL" ]; then
        echo "✅ Found schema in $FILE"
        echo "-- From: $FILE" >> "$OUTPUT_FILE"
        echo "$SQL" >> "$OUTPUT_FILE"
        echo ";" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    else
        echo "⚠️  No schema found in $FILE"
    fi
done

echo "✅ Generated $OUTPUT_FILE"
echo ""
echo "📋 Preview (first 10 lines):"
head -10 "$OUTPUT_FILE"
