import os
from prefect import task
import base64
import io
import pymupdf


@task
def decode_file(file_bytes: str) -> io.BytesIO:
    """Decode base64 string into in-memory file"""
    return io.BytesIO(base64.b64decode(file_bytes))


@task
def extract_text_from_pdf(file: io.BytesIO) -> str:
    """Extract readable text from PDF using pymupdf"""
    doc = pymupdf.open(stream=file, filetype="pdf")

    if doc.page_count == 0:
        return ""
    
    full_text = ""

    for page in doc:
        page_text = page.get_text()
        full_text += page_text + "\f"

    return full_text.strip()


@task
def extract_text_from_txt(file: io.BytesIO) -> str:
    """Read plain text file (assumed UTF-8)"""
    file.seek(0)
    return file.read().decode("utf-8").strip()


@task
def extract_text(file: io.BytesIO, file_name: str) -> str:
    """Dispatch based on file extension"""
    ALLOWED_FILE_TYPES = [".txt", ".pdf"]
    ext = os.path.splitext(file_name)[1].lower()

    if ext not in ALLOWED_FILE_TYPES:
        raise ValueError(
            f"Unsupported file type '{ext}': only {', '.join(ALLOWED_FILE_TYPES)} are allowed"
        )

    if ext == ".pdf":
        return extract_text_from_pdf(file)
    elif ext == ".txt":
        return extract_text_from_txt(file)
