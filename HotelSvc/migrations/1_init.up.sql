CREATE TABLE IF NOT EXISTS Hotels (
    ID SERIAL PRIMARY KEY,
    OwnerID INT NOT NULL,
    Name TEXT NOT NULL,
    Description TEXT
);

CREATE TABLE IF NOT EXISTS room_type (
    ID SERIAL PRIMARY KEY,
    Name TEXT NOT NULL,
    Description TEXT,
    BasePrice INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Rooms (
    ID SERIAL PRIMARY KEY,
    HotelID INT NOT NULL REFERENCES Hotels(ID) ON DELETE CASCADE,
    RoomTypeID INT NOT NULL REFERENCES room_type(ID) ON DELETE CASCADE,
    Number INT NOT NULL
);