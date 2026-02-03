// Package document defines business workflows for document processing.
// Each pipeline orchestrates a sequence of operations, such as the upload workflow:
// - Update doc metadata if needed
// - Extract metadata from documents
// - Persist metadata to the database
// - Run through OCR if needed
// - Store finalized document content in SeaweedFS object storage
package document
