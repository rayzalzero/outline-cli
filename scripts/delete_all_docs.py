#!/usr/bin/env python3
import requests

API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjA5OTc0OWNlLWNkNDMtNGJkNy1iOWI5LTEzMGMzMzkzYWFhNiIsImV4cGlyZXNBdCI6IjIwMjYtMTAtMDJUMDg6MzU6MTIuODUxWiIsInR5cGUiOiJzZXNzaW9uIiwic2VydmljZSI6ImVtYWlsIiwiaWF0IjoxNzgyOTgxMzEyfQ.xJmIhZ_yC040t4nW9qg7V44LaSZW4HzupbLsuPw2UFk"
BASE_URL = "http://localhost:3000"
COLLECTION_ID = "54ef41a8-f481-4eb7-8eb4-efd2a7407bf6"

# Get all documents
response = requests.post(
    f"{BASE_URL}/api/documents.list",
    headers={
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    },
    json={"collectionId": COLLECTION_ID, "limit": 100}
)

docs = response.json()['data']
print(f"Found {len(docs)} documents to delete\n")

# Delete each document
for doc in docs:
    print(f"Deleting: {doc['title']} ({doc['id']})")
    
    delete_response = requests.post(
        f"{BASE_URL}/api/documents.delete",
        headers={
            "Authorization": f"Bearer {API_KEY}",
            "Content-Type": "application/json"
        },
        json={"id": doc['id']}
    )
    
    if delete_response.status_code == 200:
        print(f"  ✓ Deleted")
    else:
        print(f"  ✗ Failed: {delete_response.text}")

print(f"\nDone! Deleted {len(docs)} documents")
