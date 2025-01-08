CREATE TABLE Bookings (
    ID SERIAL PRIMARY KEY,
    Status TEXT NOT NULL DEFAULT 'waiting',
    UserID INT NOT NULL,
    RoomID INT NOT NULL,
    HotelID INT NOT NULL,
    StartDate TIMESTAMP WITH TIME ZONE NOT NULL,
    EndDate TIMESTAMP WITH TIME ZONE NOT NULL,
    CreatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);