from pydantic import BaseModel


class PaymentCompletedEvent(BaseModel):

    orderId: str

    amount: float

    status: str