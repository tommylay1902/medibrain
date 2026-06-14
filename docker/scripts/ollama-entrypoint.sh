#!/bin/bash

ollama serve &

echo "Waiting for Ollama to start..."
sleep 5

ollama pull qwen3-embedding:0.6b
wait
