package com.example.demo.repository;

import org.springframework.data.repository.PagingAndSortingRepository;
import com.example.demo.entity.UUID;
import org.springframework.stereotype.Repository;

@Repository
public interface UUIDRepository
        extends PagingAndSortingRepository<UUID, Long> {
}