import uuid

from datetime import datetime

from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Float
from sqlalchemy import DateTime

from app.database.database import Base


class Notification(Base):

    __tablename__ = "notifications"

    id = Column(
        String,
        primary_key=True,
        default=lambda: str(uuid.uuid4())
    )

    order_id = Column(
        String,
        nullable=False
    )

    amount = Column(
        Float,
        nullable=False
    )

    status = Column(
        String,
        nullable=False
    )

    message = Column(
        String,
        nullable=False
    )

    created_at = Column(
        DateTime,
        default=datetime.utcnow
    )