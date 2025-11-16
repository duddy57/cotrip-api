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
  private TripRepository repository;



  @PostMapping("/create-trip")
  public ResponseEntity<TripCreateResponse> createTrip(@RequestBody TripRequestPayload payload) {
    TripModel newTripModel = new TripModel(payload);

    this.repository.save(newTripModel);

    this.participantService.registerParticipantsToEvent(payload.emails_to_invite(), newTripModel);

    return ResponseEntity.ok(
      new TripCreateResponse(newTripModel.getId())
    );

  }

  @GetMapping("/list-trips")
  public ResponseEntity<List<TripModel>> getAllTrips() {
    List<TripModel> trips = this.repository.findAll();

    if (trips.isEmpty()) {
      return ResponseEntity.<List<TripModel>>noContent().build();
    }

    return ResponseEntity.ok(trips);
  }

  @GetMapping("/{tripId}")
  public ResponseEntity<TripModel> getTripDetails(@PathVariable UUID tripId) {
    Optional<TripModel> trip = this.repository.findById(tripId);


    return trip.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
  };

  @PutMapping("/{tripId}")
  public ResponseEntity<TripModel> updateTripDetails(@PathVariable UUID tripId, @RequestBody TripRequestPayload payload) {
    Optional<TripModel> trip = this.repository.findById(tripId);

    if(trip.isEmpty()) {
      return ResponseEntity.notFound().build();
    }

    TripModel rawTripModel = trip.get();
    rawTripModel.setDestination(payload.destination());
    rawTripModel.setStartAt(LocalDateTime.parse(payload.start_at(), DateTimeFormatter.ISO_DATE_TIME));
    rawTripModel.setEndAt(LocalDateTime.parse(payload.end_at(), DateTimeFormatter.ISO_DATE_TIME));
    this.repository.save(rawTripModel);


    return ResponseEntity.ok(rawTripModel);
  };



  @GetMapping("/{tripId}/confirm")
  public ResponseEntity<TripModel> confirmTrip(@PathVariable UUID tripId) {
    Optional<TripModel> trip = this.repository.findById(tripId);

    if(trip.isEmpty()) {
      return ResponseEntity.notFound().build();
    }

    TripModel rawTripModel = trip.get();
    rawTripModel.setIsConfirmed(true);
    this.repository.save(rawTripModel);

    this.participantService.triggerConfirmationEmailToParticipants(tripId);



    return ResponseEntity.ok(rawTripModel);

  };

  @PostMapping("/{tripId}/new-invite")
  public ResponseEntity<String> inviteNewParticipants(@PathVariable UUID tripId, @RequestBody ParticipantRequestPayload payload) {
    Optional<TripModel> trip = this.repository.findById(tripId);

    if(trip.isEmpty()) {
      return ResponseEntity.notFound().build();
    }

    TripModel rawTripModel = trip.get();

    ParticipantCreatedResponse participantResponse = this.participantService.registerParticipantToEvent(payload.email(), rawTripModel.getId());

    if(rawTripModel.getIsConfirmed()) this.participantService.triggerConfirmationEmailToParticipant(payload.email());


    return ResponseEntity.ok().build();
  };

  @GetMapping("/{tripId}/participants")
  public ResponseEntity<List<ParticipantData>> getAllParticipants(@PathVariable UUID tripId) {
    List<ParticipantData> res = this.participantService.getAllParticipantsByTripId(tripId);

    return ResponseEntity.ok(res);
  }



  @PostMapping("/{tripId}/new-activity")
  public ResponseEntity<String> createNewActivity(@PathVariable UUID tripId, @RequestBody ActivityRequestPayload payload) {
    Optional<TripModel> trip = this.repository.findById(tripId);

    if(trip.isEmpty()) {
      return ResponseEntity.notFound().build();
    }

    TripModel rawTripModel = trip.get();

    ActivityResponse activityResponse = this.activitiesService.registerActivity(payload, rawTripModel.getId());

    return ResponseEntity.ok().build();
  };

  @GetMapping("/{tripId}/activities")
  public ResponseEntity<List<ActivityData>> getAllActivities(@PathVariable UUID tripId) {
    List<ActivityData> res = this.activitiesService.getAllActivitiesForTrip(tripId);
    return ResponseEntity.ok(res);
  }


  @PostMapping("/{tripId}/new-link")
  public ResponseEntity<String> createNewLink(@PathVariable UUID tripId, @RequestBody LinkRequestPayload payload) {
    Optional<TripModel> trip = this.repository.findById(tripId);

    if(trip.isEmpty()) {
      return ResponseEntity.notFound().build();
    }

    TripModel rawTripModel = trip.get();

    this.linkService.createLink(payload, rawTripModel.getId());


    return ResponseEntity.ok().build();
  }

  @GetMapping("/{tripId}/links")
  public ResponseEntity<List<LinkData>> getAllLinks(@PathVariable UUID tripId) {
    List<LinkData> res = this.linkService.getAllLinkForTrip(tripId);
    return ResponseEntity.ok(res);

  }



}
