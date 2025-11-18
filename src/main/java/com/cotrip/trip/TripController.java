package com.cotrip.trip;
import com.cotrip.activities.ActivitiesService;
import com.cotrip.activities.DTO.ActivityData;
import com.cotrip.activities.DTO.ActivityRequestPayload;
import com.cotrip.activities.DTO.ActivityResponse;
import com.cotrip.links.DTO.LinkData;
import com.cotrip.links.DTO.LinkRequestPayload;
import com.cotrip.links.DTO.LinkResponse;
import com.cotrip.links.LinkService;
import com.cotrip.participant.*;
import com.cotrip.participant.DTO.ParticipantCreatedResponse;
import com.cotrip.participant.DTO.ParticipantData;
import com.cotrip.participant.DTO.ParticipantRequestPayload;
import com.cotrip.trip.DTO.TripCreateResponse;
import com.cotrip.trip.DTO.TripGetDTO;
import com.cotrip.trip.DTO.TripRequestPayload;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;


@RestController
@RequestMapping("/trips")
public class TripController {

  @Autowired
  private ParticipantService participantService;

  @Autowired
  private ActivitiesService activitiesService;

  @Autowired
  private LinkService linkService;

  @Autowired
  private TripService tripService;

    @PostMapping("/create-trip")
    public ResponseEntity<TripCreateResponse> createTrip(@RequestBody TripRequestPayload payload) {

    var newTripModel = tripService.CreateTrip(payload);

    return ResponseEntity.ok(
      new TripCreateResponse(newTripModel.id())
    );

  }

  @GetMapping("/list-trips")
  public ResponseEntity<List<TripGetDTO>> getAllTrips() {

    var trips = tripService.getTrips();

    if (trips.isEmpty()) {
      return ResponseEntity.<List<TripGetDTO>>noContent().build();
    }

    return ResponseEntity.ok(trips);
  }

  @GetMapping("/{tripId}")
  public ResponseEntity<TripGetDTO> getTripDetails(@PathVariable UUID tripId) {

    var trip = tripService.getTripById(tripId);

    return trip.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
  };

@DeleteMapping("/{tripId}")
public ResponseEntity<String> deleteTrip(@PathVariable UUID tripId) {
    Optional<TripGetDTO> trip = this.tripService.getTripById(tripId);

    if(trip.isEmpty()) {
        return ResponseEntity.notFound().build();
    }

    tripService.DeleteTrip(tripId);

    return ResponseEntity.ok("Trip deleted successfully");
}

  @PutMapping("/{tripId}")
  public ResponseEntity<TripGetDTO> updateTripDetails(@PathVariable UUID tripId, @RequestBody TripRequestPayload payload) {

      var trip = tripService.UpdateTrip(tripId, payload);

      return trip.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());

  };

  @GetMapping("/{tripId}/confirm")
  public ResponseEntity<TripGetDTO> confirmTrip(@PathVariable UUID tripId) {

      var trip = tripService.ConfirmTrip(tripId);

      this.participantService.triggerConfirmationEmailToParticipants(tripId);

      return trip.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());

  };

    @GetMapping("/{email}/trips-by-email")
    public ResponseEntity<List<TripModel>> getTripsByEmail(@PathVariable String email) {
        List<Optional<TripGetDTO>> trips = tripService.getTripsByOwnerEmail(email);

        if (trips.isEmpty()) {
            return ResponseEntity.<List<TripModel>>noContent().build();
        }

        List<TripModel> tripModels = new ArrayList<>();

        for (Optional<TripGetDTO> tripOpt : trips) {
            tripOpt.ifPresent(tripGetDTO -> {
                TripModel tripModel = new TripModel();
                tripModel.setId(tripGetDTO.id());
                tripModel.setDestination(tripGetDTO.destination());
                tripModel.setStartAt(tripGetDTO.startAt());
                tripModel.setEndAt(tripGetDTO.endAt());
                tripModel.setIsConfirmed(tripGetDTO.isConfirmed());
                tripModels.add(tripModel);
            });
        }

        return ResponseEntity.ok(tripModels);
    }

}
