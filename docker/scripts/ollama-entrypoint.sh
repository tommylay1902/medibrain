#!/bin/bash
ollama serve &
OLLAMA_PID=$!

echo "Waiting for Ollama to start..."
until echo > /dev/tcp/localhost/11434 2>/dev/null; do
    sleep 1
done

if ! ollama list | grep -q "qwen3-embedding"; then
    echo "Pulling qwen3-embedding:0.6b..."
    ollama pull qwen3-embedding:0.6b
fi

if ! ollama list | grep -q "medgemma:4b"; then
    echo "Pulling medgemma:4b..."
    ollama pull medgemma:4b
fi

if ! ollama list | grep -q "mxbai-embed-large"; then
    echo "Pulling mxbai-embed-large..."
    ollama pull mxbai-embed-large
fi

echo "Models ready!"
wait $OLLAMA_PID
