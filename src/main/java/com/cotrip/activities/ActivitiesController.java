package com.cotrip.activities;

import com.cotrip.activities.DTO.ActivityData;
import com.cotrip.activities.DTO.ActivityRequestPayload;
import com.cotrip.activities.DTO.ActivityResponse;
import com.cotrip.trip.DTO.TripGetDTO;
import com.cotrip.trip.TripService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@RestController
@RequestMapping("/activities")
public class ActivitiesController {

    private final TripService tripService;
    private final ActivitiesService activitiesService;

    @Autowired
    public ActivitiesController(TripService tripService, ActivitiesService activitiesService) {
        this.tripService = tripService;
        this.activitiesService = activitiesService;
    }

    @PostMapping("/{tripId}/new-activity")
    public ResponseEntity<String> createNewActivity(@PathVariable UUID tripId, @RequestBody ActivityRequestPayload payload) {
        Optional<TripGetDTO> trip = tripService.getTripById(tripId);

        if(trip.isEmpty()) {
            return ResponseEntity.notFound().build();
        }
        ActivityResponse activityResponse = this.activitiesService.registerActivity(payload, trip.get().id());

        return ResponseEntity.ok().build();
    }

    @GetMapping("/{tripId}/activities")
    public ResponseEntity<List<ActivityData>> getAllActivities(@PathVariable UUID tripId) {
        List<ActivityData> res = this.activitiesService.getAllActivitiesForTrip(tripId);
        return ResponseEntity.ok(res);
    }

    @GetMapping("/{tripId}/activities/{activityId}")
    public ResponseEntity<ActivityData> getActivityById(@PathVariable UUID tripId, @PathVariable UUID activityId) {
        var res = this.activitiesService.getActivityByIdForTrip(tripId, activityId);
        return res.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @PutMapping("/{tripId}/activities/{activityId}")
    public ResponseEntity<ActivityData> updateActivity(@PathVariable UUID tripId,
                                                       @PathVariable UUID activityId,
                                                       @RequestBody ActivityRequestPayload payload) {
        var updated = this.activitiesService.updateActivity(tripId, activityId, payload);
        return updated.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{tripId}/activities/{activityId}")
    public ResponseEntity<Void> deleteActivity(@PathVariable UUID tripId, @PathVariable UUID activityId) {
        boolean deleted = this.activitiesService.deleteActivity(tripId, activityId);
        if (!deleted) return ResponseEntity.notFound().build();
        return ResponseEntity.noContent().build();
    }

}
