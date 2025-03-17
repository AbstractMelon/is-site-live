#!/bin/sh

echo "Installing dependencies..."
cd backend
go mod download
cd ..

echo "Installing frontend dependencies..."
cd frontend
pnpm install
cd ..

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

