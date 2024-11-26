# Character Search Using Go, React, Terraform and GCP
A full-stack application for managing and searching characters from the Rick and Morty API. The backend is built using Go with MongoDB as the database, and the frontend is a React application.

---

## Setup Instructions

### Prerequisites

Ensure you have the following installed on your system:
- **Go** (version 1.18 or above)
- **MongoDB** (locally running or accessible instance)
- **Node.js** (for running the frontend)
- **Git** (optional, for cloning the repository)
- **GCP** (project setup and its credentials.json file)
- **Terreform** (database connections setup)

### Backend Setup

1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd <repository_directory>/backend
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Start MongoDB locally or configure the connection string for your MongoDB instance in the backend code (default: mongosh "mongodb://localhost:27017").

4. Run the backend server:
   ```bash
   go run main.go
   ```
   By default, the backend server will run on `http://localhost:8080`.

5. Populate MongoDB with Rick and Morty character data:
   On starting the backend server, it will fetch character data from the Rick and Morty API and insert it into the MongoDB database. Duplicate data is automatically handled to prevent re-insertion.

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd <repository_directory>/frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```
   By default, the frontend will run on `http://localhost:3000`.

4. Open your browser and navigate to `http://localhost:3000` to use the application.

### Terraform Setup

1. Create a gcp-service-account.json file containing all the credentials of the gcp
2. Navigate to the terraform directory:
   ```bash
   cd <repository_directory>/gcp-terraform
   ```
3. Open terraform.tfvars file and update the credentials_file variable with your gcp-service-account.json file location
4. Run the following 3 commands in terminal:
    - Initiate - terraform init
    - Plan     - terraform plan
		- Apply    - terraform apply

---

## **API Documentation**

This API provides functionality to fetch and search for Rick and Morty character data. It connects to a MongoDB database to store character information retrieved from the Rick and Morty API.

### **Base URL**

The base URL for all endpoints is http://localhost:8080.

### **Routes**
#### *GET /search*
- **Description:** Search for characters by name.
- **Query Parameters:**
  - name (required): The name of the character to search for. The query parameter is case-insensitive.
- **Response:**
  - **200 OK:** Returns a list of characters matching the search query.
  - **400 Bad Request:** If the name query parameter is missing.
  - **500 Internal Server Error:** If there is an issue fetching or processing the data.
**Example Request:**
```bash
GET /search?name=Rick
```
**Example Response:**
```bash
[
  {
    "_id": 1,
    "name": "Rick Sanchez",
    "status": "Alive",
    "species": "Human",
    "type": "",
    "gender": "Male",
    "origin": {"name": "Earth (C-137)"},
    "location": {"name": "Earth (C-137)"},
    "image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg"
  }
]
```
---

## Architecture Overview

### Backend

- **Language**: Go
- **Frameworks**: Gorilla Mux for routing, MongoDB Go Driver for database operations
- **Database**: MongoDB
- **API**: REST API to fetch and search Rick and Morty characters.

Key features:
- Fetches character data from [Rick and Morty API](https://rickandmortyapi.com/) and stores it in MongoDB.
- Implements search functionality using MongoDB queries (case-insensitive search on the `name` field).
- Handles duplicate data gracefully using MongoDB's `upsert` mechanism.

### Frontend

- **Language**: TypeScript
- **Framework**: React
- **UI Libraries**: TailwindCSS for styling
- **HTTP Client**: Axios for API calls

Key features:
- Simple UI to search for characters by name.
- Displays character details such as name, status, species, gender, and an image.

### Workflow

1. Backend fetches and stores data from the Rick and Morty API into MongoDB.
2. React frontend communicates with the backend to search characters by name.
3. MongoDB handles the search query efficiently, and results are returned to the frontend for display.

---

## Testing Instructions

### Backend Testing

1. **API Testing**:
   - Use tools like Postman or cURL to test the `/search` endpoint:
     ```bash
     curl "http://localhost:8080/search?name=Rick"
     ```
   - Verify that the correct results are returned as a JSON array.

2. **Database Testing**:
   - Check the `rickandmorty.characters` collection in MongoDB to ensure data has been inserted correctly.
   - Use MongoDB Compass or the CLI to query the database:
     ```bash
     db.characters.find({ name: /Rick/i })
     ```

3. **Error Handling**:
   - Verify that duplicate entries do not cause crashes by restarting the backend and re-fetching data.

### Frontend Testing

1. **Search Functionality**:
   - Enter different search terms in the input field on the UI.
   - Ensure the correct results are displayed and match the backend response.

2. **Error Handling**:
   - Test with an empty search field to ensure proper validation messages.
   - Stop the backend server and verify that the frontend displays an error message when the API is unreachable.

3. **Responsive Design**:
   - Check the UI on different screen sizes to ensure proper rendering.

---

## Folder Structure

```plaintext
<repository_directory>
├── backend
│   ├── main.go           # Main backend application logic
│   ├── go.mod            # Go module dependencies
│   └── go.sum            # Go module checksums
├── frontend
│   ├── public
│   ├── src
│   │   ├── components
│   │   │   └── CharacterSearch.tsx  # Main search component
│   │   └── App.tsx
│   ├── package.json       # Node.js dependencies
│   ├── tailwind.config.js # TailwindCSS configuration
│   └── tsconfig.json      # TypeScript configuration
├── gcp-terraform
│   ├── main.tf
│   ├── outputs.tf
|   ├── terraform.tfvars
|   └── varrables.tf
└── README.md              # Project documentation
```
---

## Additional Notes

- Ensure MongoDB is running before starting the backend server.
- Adjust the MongoDB connection URI or React CORS origin as needed for production environments.
- For any issues, refer to the backend logs or browser console for debugging information.

---
## **App Demo**
This gives the brief walk through how the application will work after successfull installation on local system.

https://github.com/user-attachments/assets/dec452d4-149f-489e-ae84-9f48f6eb1180
