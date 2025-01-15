from fastapi import FastAPI, File, Form, UploadFile
from fastapi.responses import FileResponse, HTMLResponse, JSONResponse
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
    custom_name: str = Form(...)
):
    # Save the uploaded XML file
    xml_file_path = os.path.join(OUTPUT_DIR, file.filename)
    with open(xml_file_path, "wb") as xml_file:
        xml_file.write(await file.read())

    # Prepare output Excel file path
    excel_file_name = f"{custom_name}.xlsx"
    excel_file_path = os.path.join(OUTPUT_DIR, excel_file_name)

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
        return JSONResponse(content={"error": f"Error running Go backend: {e.stderr}"}, status_code=500)

    return JSONResponse(
        content={
            "message": "File processed successfully",
            "output_file": f"/files/{excel_file_name}"
        }
    )

@app.get("/files/{filename}")
async def get_file(filename: str):
    file_path = os.path.join(OUTPUT_DIR, filename)
    if not os.path.exists(file_path):
        return JSONResponse(content={"error": "File not found"}, status_code=404)
    return FileResponse(file_path, headers={"Content-Disposition": f"attachment; filename={filename}"})
