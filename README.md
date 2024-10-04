## Getting Started

### Configuration

1. Create a configuration file named `.env.local` in the root directory.
2. Define the following environment variables in the `.env.local` file:

| Key                       | Desc                                                                       |
| ------------------------- | -------------------------------------------------------------------------- |
| APP_PORT                  | Specifies the port used by the backend application                         |
| CLIENT_BASE_URL           | Base URL of the allowed frontend for communication with the backend (CORS) |
| CLOUDINARY_URL            | Complete Cloudinary URL (provided by Cloudinary service)                   |
| CLOUDINARY_CLOUD_NAME     | Cloud name on Cloudinary                                                   |
| CLOUDINARY_API_KEY        | Cloudinary API key                                                         |
| CLOUDINARY_SECRET_KEY     | Cloudinary secret key                                                      |
| CLOUDINARY_BUYCUT_FOLDER  | Folder for storing images in Cloudinary                                    |
| JWT_SECRET_KEY            | Secret key used to sign the access tokens                                  |
| JWT_ACCESS_TOKEN_DURATION | Duration of access tokens                                                  |
| POSTGRES_HOST             | Host of the PostgreSQL database                                            |
| POSTGRES_USER             | PostgreSQL username                                                        |
| POSTGRES_PASSWORD         | Password for the PostgreSQL user                                           |
| POSTGRES_DATABASE         | Name of the PostgreSQL database                                            |
| POSTGRES_PORT             | Port used by PostgreSQL                                                    |

### Setup infrastructure

- Install all dependencies

  ```
  go mod tidy
  ```

- Start a PostgreSQL database server in a Docker container:

  ```
  make dbstart
  ```

- Run the RESTful API server:

  ```
  make run
  ```
