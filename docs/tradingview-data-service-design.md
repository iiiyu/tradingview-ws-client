# Self-hosted TradingView Data Service

1. Utilize TradingView WebSocket client package
2. Provide 6 APIs:
   - Add Symbol
   - Remove Symbol
   - Get list of active symbols
   - Check symbol status
   - Get Symbol Last Price
   - Get Symbol Candlestick Data
3. Store data in SQLite database
4. Cache last price in key-value store

## Benefits

1. Unlimited API calls
2. Free of charge

## Database Design

### Database Schema for Stock Data

1. Tables Overview
   Symbols: Stores information about stock symbols.
   CandlestickData: Stores candlestick data for each stock symbol and time period.
   TimePeriods: Stores different time periods (e.g., 1 minute, 5 minutes, 1 hour).
2. Table Definitions
   Symbols Table

```sql
CREATE TABLE Symbols (
id INTEGER PRIMARY KEY AUTOINCREMENT,
symbol TEXT UNIQUE NOT NULL,
name TEXT NOT NULL
);
```

id: Primary key for reference.
symbol: Stock ticker symbol (e.g., AAPL for Apple).
name: Full name of the stock or company.
TimePeriods Table

```sql
CREATE TABLE TimePeriods (
id INTEGER PRIMARY KEY AUTOINCREMENT,
period TEXT UNIQUE NOT NULL
);
```

id: Primary key for reference.
period: Time period name (e.g., "1min", "5min", "1hour").
CandlestickData Table

```sql
CREATE TABLE CandlestickData (
id INTEGER PRIMARY KEY AUTOINCREMENT,
symbol_id INTEGER NOT NULL,
period_id INTEGER NOT NULL,
open REAL NOT NULL,
high REAL NOT NULL,
low REAL NOT NULL,
close REAL NOT NULL,
volume REAL NOT NULL,
timestamp DATETIME NOT NULL,
FOREIGN KEY (symbol_id) REFERENCES Symbols(id),
FOREIGN KEY (period_id) REFERENCES TimePeriods(id),
UNIQUE (symbol_id, period_id, timestamp)
);
```

symbol_id: Foreign key linking to the Symbols table.
period_id: Foreign key linking to the TimePeriods table.
open: Opening price.
high: Highest price.
low: Lowest price.
close: Closing price.
volume: Trading volume.
timestamp: Date and time of the candlestick.

### Considerations

Indexes: Create indexes on frequently queried columns (symbol_id, period_id, timestamp) for faster lookups.

```sql
CREATE INDEX idx_candlestick_symbol_period_timestamp ON CandlestickData(symbol_id, period_id, timestamp);
```

Data Retention: Implement a cleanup script to remove old data if storage becomes an issue.
Scalability: For larger datasets, consider partitioning by symbol or time period in your application logic.
