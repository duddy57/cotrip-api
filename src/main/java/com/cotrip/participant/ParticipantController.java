package com.cotrip.participant;


import com.cotrip.participant.DTO.ParticipantRequestPayload;
import com.cotrip.participant.DTO.ParticipantData;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Optional;
import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/participants")
public class ParticipantController {

  @Autowired
  private ParticipantService participantService;


  @PostMapping("/{id}/confirm")
  public ResponseEntity<String> confirmParticipation(@PathVariable UUID id, @RequestBody ParticipantRequestPayload payload) {
  boolean ok = this.participantService.confirmParticipant(id, payload);
  if (!ok) return ResponseEntity.notFound().build();
  return ResponseEntity.ok("Participation confirmed for participant ID: " + id);
  };

  @GetMapping("/{tripId}")
  public ResponseEntity<List<ParticipantData>> listByTrip(@PathVariable UUID tripId) {
    return ResponseEntity.ok(this.participantService.getAllParticipantsByTripId(tripId));
  }

  @GetMapping("/by-id/{id}")
  public ResponseEntity<ParticipantData> getById(@PathVariable UUID id) {
    return this.participantService.getParticipantById(id)
            .map(ResponseEntity::ok)
            .orElseGet(() -> ResponseEntity.notFound().build());
  }

  @PutMapping("/{id}")
  public ResponseEntity<String> updateDetails(@PathVariable UUID id, @RequestBody ParticipantRequestPayload payload) {
    boolean ok = this.participantService.updateParticipantDetails(id, payload);
    return ok ? ResponseEntity.ok("Participant Updated") : ResponseEntity.notFound().build();
  }

  @DeleteMapping("/{id}")
  public ResponseEntity<String> delete(@PathVariable UUID id) {
    boolean ok = this.participantService.deleteParticipantsById(id);
    return ok ? ResponseEntity.ok("Participant Deleted"): ResponseEntity.notFound().build();
  }

}
