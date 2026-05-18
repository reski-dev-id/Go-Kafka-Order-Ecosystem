import json

from aiokafka import AIOKafkaConsumer

from app.dto.payment_completed_event import (
    PaymentCompletedEvent,
)

from app.database.database import SessionLocal

from app.entity.notification import Notification


async def consume_payment_completed():

    consumer = AIOKafkaConsumer(
        "payment.completed",
        bootstrap_servers="localhost:9094",
        group_id="notification-group",
        auto_offset_reset="earliest"
    )

    await consumer.start()

    print("notification consumer started")

    try:

        async for message in consumer:

            data = json.loads(
                message.value.decode("utf-8")
            )

            event = PaymentCompletedEvent(**data)

            db = SessionLocal()

            notification = Notification(
                order_id=event.orderId,
                amount=event.amount,
                status=event.status,
                message=f"payment success for order {event.orderId}"
            )

            db.add(notification)

            db.commit()

            db.close()

            print(
                f"notification created for order: "
                f"{event.orderId}"
            )

    finally:

        await consumer.stop()