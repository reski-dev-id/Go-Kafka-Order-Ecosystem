package com.reski.payment.dto;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class PaymentCompletedEvent {

    private String orderId;

    private Double amount;

    private String status;
}