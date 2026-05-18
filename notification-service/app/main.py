import asyncio

from fastapi import FastAPI

from app.database.database import engine
from app.database.database import Base

from app.consumer.payment_completed_consumer import (
    consume_payment_completed,
)

app = FastAPI()

Base.metadata.create_all(bind=engine)


@app.on_event("startup")
async def startup_event():

    asyncio.create_task(
        consume_payment_completed()
    )


@app.get("/health")
async def health():

    return {
        "success": True,
        "message": "notification-service running"
    }