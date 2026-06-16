from fastapi import FastAPI, UploadFile, File, HTTPException
from uuid import uuid4
from pathlib import Path
from langchain_community.document_loaders import PyPDFLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter
import magic
import os
import tempfile

app = FastAPI()

MAX_SIZE = 10 * 1024 * 1024  # 10MB


@app.post("/api/v1/document/split")
async def ingest_document(fileInput: UploadFile = File(...)):
    contents = await fileInput.read()

    if len(contents) > MAX_SIZE:
        raise HTTPException(status_code=400, detail="File too large")

    mime = magic.from_buffer(contents, mime=True)
    if mime != "application/pdf":
        raise HTTPException(status_code=400, detail="Only PDFs allowed")

    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf", dir="/tmp") as tmp:
        tmp.write(contents)
        tmp_path = tmp.name

    try:
        loader = PyPDFLoader(tmp_path)
        pages = loader.load()
        splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
            is_separator_regex=False,
        )
        splits = splitter.split_documents(pages)
        print(
            f"Split into {len(splits)} chunks. Indexing now (watch for ✓ batch lines)..."
        )
    finally:
        os.unlink(tmp_path)

    return {
        "filename": fileInput.filename,
        "total_pages": len(pages),
        "chunks": [
            {"page": s.metadata["page"] + 1, "text": s.page_content}
            for i, s in enumerate(splits)
        ],
    }
