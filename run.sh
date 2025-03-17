#!/bin/sh

echo "Building binaries..."
cd backend
go build -o server cmd/server/main.go
cd ..

echo "Starting services..."
cd backend
./server &
cd ..

echo "Starting frontend development server..."
cd frontend
pnpm run dev
