#!/bin/bash

curl -i -X POST http://localhost:3000/api/v1/pages/new -H "Content-Type: application/json" -d '{"page_title":"value", "content":"testing"}'
