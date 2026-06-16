#!/bin/bash
ollama serve &
OLLAMA_PID=$!

echo "Waiting for Ollama to start..."
until echo > /dev/tcp/localhost/11434 2>/dev/null; do
    sleep 1
done

# Only pull if not already downloaded
if ! ollama list | grep -q "gemma3:4b"; then
    echo "Pulling gemma3:4b..."
    ollama pull gemma3:4b
fi

if ! ollama list | grep -q "qwen3-embedding"; then
    echo "Pulling qwen3-embedding:0.6b..."
    ollama pull qwen3-embedding:0.6b
fi

echo "Models ready!"
wait $OLLAMA_PID
