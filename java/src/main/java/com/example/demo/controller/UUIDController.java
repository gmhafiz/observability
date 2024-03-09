package com.example.demo.controller;

import com.example.demo.entity.UUID;
import com.example.demo.service.UUIDService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.logging.Logger;

@RestController
@RequestMapping("/uuid")
public class UUIDController
{
    @Autowired
    UUIDService service;

    @GetMapping
    public ResponseEntity<List<UUID>> getAllUUID(
            @RequestParam(defaultValue = "0") Integer pageNo,
            @RequestParam(defaultValue = "10") Integer pageSize,
            @RequestParam(defaultValue = "id") String sortBy)
    {
        System.out.println("entering getAllUUID controller");

        List<UUID> list;
        try {
            list = service.getUUIDs(pageNo, pageSize, sortBy);
        } catch (Exception e) {
            return new ResponseEntity<>(null, new HttpHeaders(), HttpStatus.INTERNAL_SERVER_ERROR);
        }

        return new ResponseEntity<>(list, new HttpHeaders(), HttpStatus.OK);
    }
}