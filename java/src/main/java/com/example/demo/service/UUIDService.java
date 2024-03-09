package com.example.demo.service;

import com.example.demo.repository.UUIDRepository;
import com.example.demo.entity.UUID;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

@Service
public class UUIDService
{

    final
    UUIDRepository repository;
    private final Random rand;


    public UUIDService(UUIDRepository repository) {
        this.repository = repository;

        this.rand = new Random();
    }

    public List<UUID> getUUIDs(Integer pageNo, Integer pageSize, String sortBy) throws Exception {

        float val = this.rand.nextFloat();
        if (val < 0.05) {
            System.out.println("random database fail");
            throw new Exception("random database fail");
        }

        Pageable paging = PageRequest.of(pageNo, pageSize, Sort.by(sortBy).descending());

        Page<UUID> pagedResult = repository.findAll(paging);

        if(pagedResult.hasContent()) {
            return pagedResult.getContent();
        } else {
            return new ArrayList<UUID>();
        }
    }
}