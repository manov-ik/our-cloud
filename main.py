from fastapi import FastAPI, File, UploadFile, Request
from fastapi.responses import HTMLResponse, FileResponse, JSONResponse
from fastapi.staticfiles import StaticFiles
import os
from pathlib import Path

app = FastAPI()

UPLOAD_DIR = Path("shared")
UPLOAD_DIR.mkdir(exist_ok=True)

app.mount("/static", StaticFiles(directory="static"), name="static")

@app.get("/", response_class=HTMLResponse)
async def serve_html():
    return (UPLOAD_DIR.parent / "static/index.html").read_text(encoding="utf-8")

@app.post("/upload")
async def upload_file(file: UploadFile = File(...)):
    file_path = UPLOAD_DIR / file.filename
    with open(file_path, "wb") as f:
        content = await file.read()
        f.write(content)
    return {"message": "Uploaded!"}

@app.get("/files")
async def list_files():
    files = [f.name for f in UPLOAD_DIR.iterdir() if f.is_file()]
    return JSONResponse(content=files)

@app.get("/download")
async def download_file(name: str):
    file_path = UPLOAD_DIR / name
    if file_path.exists():
        return FileResponse(file_path, filename=name)
    return JSONResponse(status_code=404, content={"error": "File not found"})
