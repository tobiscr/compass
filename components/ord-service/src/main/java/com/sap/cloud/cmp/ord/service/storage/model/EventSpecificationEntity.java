package com.sap.cloud.cmp.ord.service.storage.model;

import org.eclipse.persistence.annotations.Convert;
import org.eclipse.persistence.annotations.TypeConverter;

import javax.persistence.*;
import java.sql.Clob;
import java.util.UUID;

/**
 * Once we introduce ORD Aggregator EventDefinitions should be fetched as Specifications.
 * Therefore Compass should support multiple specs per Event.
 * In order to achive that we will create a new specifications table which will unify api/event specifications.
 */
@Entity(name = "eventSpecification")
@Table(name = "event_api_definitions")
public class EventSpecificationEntity {
    @javax.persistence.Id
    @Column(name = "id")
    @Convert("uuidConverter")
    @TypeConverter(name = "uuidConverter", dataType = Object.class, objectType = UUID.class)
    private UUID eventDefinitionId;

    @Column(name = "spec_data", length = Integer.MAX_VALUE)
    @Lob
    private Clob specData;
}
