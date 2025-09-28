FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    ffmpeg \
    build-essential \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy and install dependencies
COPY requirements.txt .

# Upgrade pip and install dependencies
RUN pip install --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Copy cookies file
COPY cookies.txt /app/

# Create non-root user and fix permissions
RUN useradd --create-home appuser && \
    chown -R appuser:appuser /app
USER appuser

# Expose gRPC port
EXPOSE 50051

# Start the service
CMD ["python", "main.py"]