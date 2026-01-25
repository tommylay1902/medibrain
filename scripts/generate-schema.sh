#!/bin/bash
# scripts/extract-schema.sh

# Default output file name
OUTPUT_FILE="../internal/database/migrations/schemas.sql"

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [-o output_file.sql]"
            echo "Options:"
            echo "  -o, --output    Specify output file name (default: schemas.sql)"
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

# Clear the output file
> "../internal/database/$OUTPUT_FILE"

echo "🔍 Searching for schemas..."
echo "📁 Output file: $OUTPUT_FILE"

find ../internal/api/domain -name "model.go" | while read FILE; do
    echo "📄 Processing: $FILE"
    
    # Read the entire file
    CONTENT=$(cat "$FILE")
    
    # Find everything between backticks after "var.*Schema.*="
    # Using Perl for multi-line matching
    SQL=$(echo "$CONTENT" | perl -ne '
        if (/var.*[Ss]chema.*=\s*`(.*)/) {
            $in_sql = 1;
            $sql = $1;
            if ($sql =~ /`/) {
                # Backtick on same line
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
