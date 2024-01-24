# Unit of work showcase

1. Clone the repository
2. Start the database
    ```shell
    docker-compose up -d
    ```

3. Run the service with

    ```shell
    go run github.com/bread-amammay/unitofwork/cmd/unitofwork
    ```

4. Run evans to interact with the service
   To install Evans
    ```shell
    brew tap ktr0731/evans
    brew install evans
    ```
   To run Evans
    ```shell
    evans repl --reflection --host=localhost -p 8080 --header=X-User-Id=$(uuidgen),X-User-Name=$HOST,X-First-Name=$USER,X-Last-Name=$USER
    ```
