package com.reski.payment.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class OrderCreatedEvent {

    @JsonProperty("ID")
    private String id;

    @JsonProperty("CustomerName")
    private String customerName;

    @JsonProperty("ProductName")
    private String productName;

    @JsonProperty("Quantity")
    private Integer quantity;

    @JsonProperty("Amount")
    private Double amount;

    @JsonProperty("Status")
    private String status;

    @JsonProperty("CreatedAt")
    private String createdAt;

    @JsonProperty("UpdatedAt")
    private String updatedAt;
}