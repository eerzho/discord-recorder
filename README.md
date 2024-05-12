Quick Start:

1. Set Up Environment Variables <br>
    * Create an Environment File: Duplicate the provided `.env.example` file and rename it to `.env`. Fill out the necessary details. <br>
2. Configure MinIO <br>
    * Run the MinIO using Docker compose:
        * Execute the command: `docker compose up -d --build` <br>
    * Set Up MinIO: <br>
      * Access the MinIO web interface by navigating to `http://localhost:9090`. <br>
      * Create a Bucket: In the MinIO interface, create a new bucket and note the bucket name. <br>
      * Configure Access Keys: Create an access key and a secret key. <br>
    * Update .env File: <br>
      * Save the bucket name under the key `MINIO_BUCKET`. <br>
      * Save the access key and secret key under the keys `MINIO_ACCESS_KEY` and `MINIO_SECRET_KEY`, respectively. <br>
3. Set Up Discord Bot <br>
    * Create and Configure Discord Bot: <br>
      * Go to the Discord Developer Portal and create a new bot. <br>
      * Save the bot token provided by Discord in your .env file under the key `DISCORD_TOKEN`. <br>
      * Add the Bot to Your Guild: Generate an invite link with the necessary permissions and add the bot to your Discord guild. <br>
4. Launch the Application <br>
    * Start the Bot: <br>
      * Run the command: `go run ./cmd/app` <br>
