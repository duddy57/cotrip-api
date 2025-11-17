package com.cotrip.activities;

import com.cotrip.activities.DTO.ActivityData;
import com.cotrip.activities.DTO.ActivityRequestPayload;
import com.cotrip.activities.DTO.ActivityResponse;
import com.cotrip.trip.TripRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.UUID;

@Service
public class ActivitiesService {

  @Autowired
  private ActivitiesRepository repository;

  @Autowired
  private TripRepository tripRepository;

  @CacheEvict(key = "#tripId" , value = "Activity")
  public ActivityResponse registerActivity(ActivityRequestPayload payload, UUID tripId) {

      var tripModel = tripRepository.findById(tripId).orElse(null);
    ActivitiesModel newActivity = new ActivitiesModel(payload.title(), payload.occurs_at(), tripModel);

    this
      .repository
      .save(newActivity);

    return new ActivityResponse(newActivity.getId());
  }

  @Cacheable("Activity")
  public List<ActivityData> getAllActivitiesForTrip(UUID trip_id) {
    return this
      .repository
      .findByTripId(trip_id)
      .stream()
      .map(activity -> new ActivityData(activity.getId(), activity.getTitle(), activity.getOccursAt()))
      .toList();
  }

  public void deleteActivityById(UUID activityId) {
    this.repository.deleteById(activityId);
  }

  public void updateActivitiesById(UUID activityId, ActivityRequestPayload payload) {
    var activityOpt = this.repository.findById(activityId);
    if (activityOpt.isPresent()) {
      var activity = activityOpt.get();
      activity.setTitle(payload.title());
      activity.setOccursAt(LocalDateTime.parse(payload.occurs_at(), DateTimeFormatter.ISO_LOCAL_DATE));
      this.repository.save(activity);
    }
  }
}
