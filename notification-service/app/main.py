import asyncio
import logging

from contextlib import asynccontextmanager

from fastapi import FastAPI

from app.database.database import engine
from app.database.database import Base

from app.consumer.payment_completed_consumer import (
    consume_payment_completed,
)

logging.basicConfig(
    level=logging.INFO,
)

logger = logging.getLogger(
    "notification-service"
)

consumer_task = None

Base.metadata.create_all(bind=engine)


@asynccontextmanager
async def lifespan(app: FastAPI):

    global consumer_task

    logger.info(
        "notification-service starting"
    )

    consumer_task = asyncio.create_task(
        consume_payment_completed()
    )

    yield

    logger.info(
        "shutdown signal received"
    )

    if consumer_task:

        logger.info(
            "stopping kafka consumer task"
        )

        consumer_task.cancel()

        try:
            await consumer_task

        except asyncio.CancelledError:

            logger.info(
                "consumer task cancelled"
            )

    logger.info(
        "notification-service shutdown complete"
    )


app = FastAPI(
    lifespan=lifespan
)


@app.get("/health")
async def health():

    return {
        "success": True,
        "message": "notification-service running"
    }