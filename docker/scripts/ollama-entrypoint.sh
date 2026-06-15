#!/bin/bash
ollama serve &
OLLAMA_PID=$!

echo "Waiting for Ollama to start..."
until echo > /dev/tcp/localhost/11434 2>/dev/null; do
    sleep 1
done

# Only pull if not already downloaded
if ! ollama list | grep -q "llama3.2:1b"; then
    echo "Pulling llama3.2:1bb..."
    ollama pull llama3.2:1b
fi

if ! ollama list | grep -q "qwen3-embedding"; then
    echo "Pulling qwen3-embedding:0.6b..."
    ollama pull qwen3-embedding:0.6b
fi

echo "Models ready!"
wait $OLLAMA_PID
