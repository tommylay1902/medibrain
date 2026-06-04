#!/bin/bash

ollama serve &

echo "Waiting for Ollama to start..."
sleep 5

ollama pull all-minilm:l6-v2

wait
