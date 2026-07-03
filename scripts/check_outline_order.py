#!/usr/bin/env python3
import requests
from datetime import datetime

API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjA5OTc0OWNlLWNkNDMtNGJkNy1iOWI5LTEzMGMzMzkzYWFhNiIsImV4cGlyZXNBdCI6IjIwMjYtMTAtMDJUMDg6MzU6MTIuODUxWiIsInR5cGUiOiJzZXNzaW9uIiwic2VydmljZSI6ImVtYWlsIiwiaWF0IjoxNzgyOTgxMzEyfQ.xJmIhZ_yC040t4nW9qg7V44LaSZW4HzupbLsuPw2UFk"
BASE_URL = "http://localhost:3000"
COLLECTION_ID = "54ef41a8-f481-4eb7-8eb4-efd2a7407bf6"

response = requests.post(
    f"{BASE_URL}/api/documents.list",
    headers={
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json"
    },
    json={"collectionId": COLLECTION_ID, "limit": 50}
)

data = response.json()
docs = data['data']

doc_map = {doc['id']: doc for doc in docs}
children_map = {}
root_docs = []

for doc in docs:
    parent_id = doc.get('parentDocumentId')
    if parent_id:
        if parent_id not in children_map:
            children_map[parent_id] = []
        children_map[parent_id].append(doc)
    else:
        root_docs.append(doc)

def print_tree(doc, prefix="", is_last=True):
    connector = "└── " if is_last else "├── "
    created = datetime.fromisoformat(doc['createdAt'].replace('Z', '+00:00'))
    print(f"{prefix}{connector}{doc['title']} ({created.strftime('%H:%M:%S')})")
    
    children = children_map.get(doc['id'], [])
    children_sorted = sorted(children, key=lambda x: x['createdAt'], reverse=True)
    
    new_prefix = prefix + ("    " if is_last else "│   ")
    for i, child in enumerate(children_sorted):
        print_tree(child, new_prefix, i == len(children_sorted) - 1)

print("📁 Outline Collection Tree (newest first, with hierarchy):\n")

root_sorted = sorted(root_docs, key=lambda x: x['createdAt'], reverse=True)
for i, doc in enumerate(root_sorted):
    print_tree(doc, "", i == len(root_sorted) - 1)
