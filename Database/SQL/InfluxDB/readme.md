# Setting up InfluxDB on macOS

## Step 1: Install InfluxDB
1. Open your terminal.
2. Install InfluxDB using Homebrew:
   ```bash
   brew update
   brew install influxdb
   ```

## Step 2: Start the InfluxDB Service
1. Start the InfluxDB service:
   ```bash
   influxd
   ```
2. By default, InfluxDB will run on `http://localhost:8086`.

## Step 3: Verify Installation
1. Open another terminal window.
2. Use the `influx` CLI to connect to the InfluxDB instance:
   ```bash
   influx
   ```
3. You should see the InfluxDB shell prompt (`>`).

## Step 4: Create a Database
1. In the InfluxDB shell, create a new database:
   ```sql
   CREATE DATABASE my_database;
   ```
2. Verify the database was created:
   ```sql
   SHOW DATABASES;
   ```

## Step 5: Write and Query Data
1. Write a data point into the database:
   ```bash
   curl -i -XPOST 'http://localhost:8086/write?db=my_database' --data-binary 'weather,location=us-midwest temperature=82 1465839830100400200'
   ```
2. Query the data:
   ```sql
   SELECT * FROM weather;
   ```

## Step 6: Stop the InfluxDB Service
1. To stop the InfluxDB service, press `Ctrl+C` in the terminal running `influxd`.

For more details, refer to the [official InfluxDB documentation](https://docs.influxdata.com/influxdb/).
