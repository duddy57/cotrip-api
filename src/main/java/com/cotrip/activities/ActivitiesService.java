package com.cotrip.activities;

import com.cotrip.activities.DTO.ActivityData;
import com.cotrip.activities.DTO.ActivityRequestPayload;
import com.cotrip.activities.DTO.ActivityResponse;
import com.cotrip.trip.TripModel;
import com.cotrip.trip.TripRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.Caching;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
public class ActivitiesService {

  @Autowired
  private ActivitiesRepository repository;

  @Autowired
  private TripRepository tripRepository;

  @Caching(evict = {
    @CacheEvict(key = "#tripId", value = "Activity"),
    @CacheEvict(value = "ActivityList", allEntries = true)
  })
  public ActivityResponse registerActivity(ActivityRequestPayload payload, UUID tripId) {

      var tripModel = tripRepository.findById(tripId).orElse(null);
    ActivitiesModel newActivity = new ActivitiesModel(payload.title(), payload.occurs_at(), tripModel);

    this
      .repository
      .save(newActivity);

    return new ActivityResponse(newActivity.getId());
  }

  @Cacheable(value = "ActivityList", key = "'all'")
  public List<ActivityData> getAllActivitiesForTrip(UUID trip_id) {
    return this
      .repository
      .findByTripId(trip_id)
      .stream()
      .map(activity -> new ActivityData(activity.getId(), activity.getTitle(), activity.getOccursAt()))
      .toList();
  }

  @Cacheable("Activity")
  public Optional<ActivityData> getActivityById(UUID activityId) {
    return this.repository.findById(activityId)
      .map(activity -> new ActivityData(activity.getId(), activity.getTitle(), activity.getOccursAt()));
  }

  @Cacheable(value = "ActivityItem", key = "#tripId.toString() + ':' + #activityId.toString()")
  public Optional<ActivityData> getActivityByIdForTrip(UUID tripId, UUID activityId) {
    Optional<ActivitiesModel> optional = this.repository.findById(activityId);
    if (optional.isEmpty()) return Optional.empty();
    ActivitiesModel entity = optional.get();
    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return Optional.empty();
    }
    return Optional.of(new ActivityData(entity.getId(), entity.getTitle(), entity.getOccursAt()));
  }

  @Caching(evict = {
    @CacheEvict(key = "#tripId", value = "Activity"),
    @CacheEvict(value = "ActivityList", allEntries = true),
    @CacheEvict(key = "#tripId.toString() + ':' + #activityId.toString()", value = "ActivityItem")
  })
  public Optional<ActivityData> updateActivity(UUID tripId, UUID activityId, ActivityRequestPayload payload) {
    Optional<ActivitiesModel> optional = this.repository.findById(activityId);

    if (optional.isEmpty()) return Optional.empty();

    ActivitiesModel entity = optional.get();

    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return Optional.empty();
    }

    if (payload.title() != null) entity.setTitle(payload.title());
    if (payload.occurs_at() != null) {
      entity.setOccursAt(LocalDateTime.parse(payload.occurs_at(), DateTimeFormatter.ISO_DATE_TIME));
    }

    this.repository.save(entity);

    return Optional.of(new ActivityData(entity.getId(), entity.getTitle(), entity.getOccursAt()));
  }

  @Caching(evict = {
    @CacheEvict(key = "#tripId", value = "Activity"),
    @CacheEvict(value = "ActivityList", allEntries = true),
    @CacheEvict(key = "#tripId.toString() + ':' + #activityId.toString()", value = "ActivityItem")
  })
  public boolean deleteActivity(UUID tripId, UUID activityId) {
    Optional<ActivitiesModel> optional = this.repository.findById(activityId);
    if (optional.isEmpty()) return false;

    ActivitiesModel entity = optional.get();
    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return false;
    }

    this.repository.delete(entity);
    return true;
  }
}
