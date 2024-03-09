package com.example.demo.entity;

import jakarta.persistence.*;


@Entity
@Table(name = "uuid")
public class UUID {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer id;


    private java.util.UUID uuid;


    public UUID(String uuid) {
        this.uuid = java.util.UUID.fromString(uuid);
    }

    public java.util.UUID getUuid() {
        return uuid;
    }

    public void setUuid(java.util.UUID uuid) {
        this.uuid = uuid;
    }

    public UUID() {}

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }
}
