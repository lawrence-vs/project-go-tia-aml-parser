function downloadFile(outputFile) {
    const downloadUrl = `http://127.0.0.1:8000${outputFile}`;
    const a = document.createElement('a');
    a.href = downloadUrl;
    a.download = outputFile.split('/').pop(); // Extracts the filename
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
}

async function uploadFile(event) {
    event.preventDefault();
    const formData = new FormData(event.target);

    try {
        const response = await fetch("/upload/", {
            method: "POST",
            body: formData
        });

        const data = await response.json();
        if (response.ok) {
            // Show popup with the success message
            alert(`Success: ${data.message}`);

            // Create a download link
            const downloadLink = document.createElement("a");
            downloadLink.href = data.output_file;
            downloadLink.download = true; // Suggest downloading
            downloadLink.textContent = "Click here to download the file";
            document.getElementById("download-container").appendChild(downloadLink);

            // Redirect to the main page after some time
            setTimeout(() => {
                window.location.href = "/";
            }, 5000); // 5 seconds
        } else {
            alert(`Error: ${data.error || "An error occurred."}`);
        }
    } catch (error) {
        console.error("Upload failed:", error);
        alert("An error occurred while uploading the file.");
    }
}