# Self-hosted TradingView Data Service

1. Utilize TradingView WebSocket client package
2. Provide 6 APIs:
   - Add Symbol(DB)
   - Remove Symbol(DB)
   - Get list of active symbols(DB)
   - Check symbol status(is active or not)(DB)
   - Get Symbol Last Price(from key-value store)
   - Get Symbol Candlestick Data(from PostgreSQL database)
3. Store data in PostgreSQL database
4. Cache last price in key-value store

## Benefits

1. Unlimited API calls
2. Free of charge

## Data Design

1. Tables Overview

   - Active Symbols: Maintains mapping between TradingView session IDs and symbols, where each symbol is identified by exchange and ticker. Each session ID maps to a unique combination of exchange, symbol, and time period.
   - Candlestick Data: Stores OHLCV (Open, High, Low, Close, Volume) data for each symbol across multiple timeframes. Each symbol maintains data for four time periods (10 seconds, 1 minute, 5 minutes, 1 day). The combination of exchange, symbol, and time period serves as a unique identifier.

2. Key-Value Store Overview
   - Real-time Market Data: Maintains the latest market data for each symbol (identified by exchange and ticker), including:
     - Last update timestamp
     - Opening price
     - High price
     - Low price
     - Closing price
     - Trading volume

### SQL Schema

```sql
-- Create enum for timeframes
-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum for timeframes (unchanged)
CREATE TYPE timeframe AS ENUM ("10S", "1", "5", "1D");

-- Renamed from trading_sessions to active_sessions
CREATE TABLE active_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id VARCHAR(50) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    timeframe timeframe NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,  -- Added enabled field
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id)
);

-- Update candles table to use UUID
CREATE TABLE candles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    exchange VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    timeframe timeframe NOT NULL,
    timestamp BIGINT NOT NULL,  -- Unix timestamp from TradingView
    open DECIMAL(20, 8) NOT NULL,
    high DECIMAL(20, 8) NOT NULL,
    low DECIMAL(20, 8) NOT NULL,
    close DECIMAL(20, 8) NOT NULL,
    volume DECIMAL(20, 8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Create a unique constraint to prevent duplicate entries
    UNIQUE(exchange, symbol, timeframe, timestamp)
);

-- Rest of the schema remains the same
CREATE INDEX idx_candles_symbol ON candles(exchange, symbol);
CREATE INDEX idx_candles_timestamp ON candles(timestamp);
CREATE INDEX idx_candles_lookup ON candles(exchange, symbol, timeframe, timestamp);

```
