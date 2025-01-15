from fastapi import FastAPI, File, Form, UploadFile
from fastapi.responses import FileResponse, HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from fastapi.requests import Request
import os
import subprocess

app = FastAPI()

# Mount static and template directories
app.mount("/static", StaticFiles(directory="static"), name="static")
templates = Jinja2Templates(directory="templates")

# Path to the Go application (located in the backend directory)
GO_MAIN_FILE = os.path.abspath("../backend/main.go")
OUTPUT_DIR = os.path.abspath("../backend/")

@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})

@app.post("/upload/")
async def upload_file(
    file: UploadFile = File(...),
    custom_name: str = Form(...),
):
    # Save the uploaded XML file to the backend directory
    xml_file_path = os.path.join(OUTPUT_DIR, file.filename)
    with open(xml_file_path, "wb") as xml_file:
        xml_file.write(await file.read())

    # Prepare output Excel file path in the backend directory
    excel_file_path = os.path.join(OUTPUT_DIR, f"{custom_name}.xlsx")

    # Run the Go application to process the XML file
    try:
        result = subprocess.run(
            ["go", "run", GO_MAIN_FILE],
            cwd=OUTPUT_DIR,
            capture_output=True,
            text=True,
            check=True
        )
    except subprocess.CalledProcessError as e:
        return {"error": f"Error running Go backend: {e.stderr}"}

    # Return a response with the generated file's path
    return {
        "message": "File processed successfully",
        "output_file": f"/files/{custom_name}.xlsx"
    }

@app.get("/files/{filename}")
async def get_file(filename: str):
    file_path = os.path.join(OUTPUT_DIR, filename)
    if not os.path.exists(file_path):
        return {"error": "File not found"}
    return FileResponse(file_path)
