#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo " - Health check "
curl -s "$BASE_URL/health"
echo

echo " - Admin login "
ADMIN_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bitspilani.ac.in","password":"admin123"}')
echo "$ADMIN_LOGIN"
ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | jq -r '.token')

echo

echo " - Register faculty "
curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Dr Smith","email":"faculty@university.edu","password":"faculty123","role":"faculty","dept":"CS"}'
echo

echo " - Faculty login "
FACULTY_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"faculty@university.edu","password":"faculty123"}')
echo "$FACULTY_LOGIN"
FACULTY_TOKEN=$(echo "$FACULTY_LOGIN" | jq -r '.token')

echo

echo " - Register student "
curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"student@university.edu","password":"student123","role":"student","dept":"CS"}'
echo

echo " - Student login "
STUDENT_LOGIN=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"student@university.edu","password":"student123"}')
STUDENT_TOKEN=$(echo "$STUDENT_LOGIN" | jq -r '.token')
STUDENT_ID=$(echo "$STUDENT_LOGIN" | jq -r '.user.id')

echo " - Student applies for leave "
APPLY_LEAVE=$(curl -s -X POST "$BASE_URL/leaves/apply" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -d "{\"student_id\":$STUDENT_ID,\"leave_type\":\"Medical\",\"reason\":\"Fever\",\"start_date\":\"2025-10-21\",\"end_date\":\"2025-10-23\"}")
echo "$APPLY_LEAVE"
LEAVE_ID=$(echo "$APPLY_LEAVE" | jq -r '.id')
echo

echo " - Faculty approves leave "
curl -s -X PUT "$BASE_URL/leaves/$LEAVE_ID/approve" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $FACULTY_TOKEN" \
  -d '{"status":"approved"}'
echo

echo " - Mark attendance "
curl -s -X POST "$BASE_URL/attendance/mark" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $FACULTY_TOKEN" \
  -d "{\"student_id\":$STUDENT_ID,\"date\":\"2025-10-20\",\"present\":true}"
echo

echo " - Attendance stats "
curl -s "$BASE_URL/attendance/stats/$STUDENT_ID" -H "Authorization: Bearer $FACULTY_TOKEN" | jq .
echo

echo "Tests completed."
