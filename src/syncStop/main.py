from fastapi import FastAPI
from routes import songs
import uvicorn
from config import config

app = FastAPI(
    title="SyncStop",
    description="SyncStop is a service that allows you to access your music library",
    version="1.0.0"
)

app.include_router(songs.router, prefix="/api/v1")

@app.get("/")
def read_root():
    return {"message": "Hello World"}

if __name__ == "__main__":
    uvicorn.run(app, host=config.HOST, port=config.PORT)