package com.cotrip.participant;


import com.cotrip.participant.DTO.ParticipantCreatedResponse;
import com.cotrip.participant.DTO.ParticipantData;
import com.cotrip.participant.DTO.ParticipantRequestPayload;
import com.cotrip.trip.DTO.TripGetDTO;
import com.cotrip.trip.TripModel;
import com.cotrip.trip.TripRepository;
import com.cotrip.trip.TripService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.Caching;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.UUID;

@Service
public class ParticipantService {

  @Autowired
  private ParticipantRepository repository;

  @Autowired
  private TripRepository tripRepository;

  @CacheEvict(value = "Participant", key = "#trip.id")
  public void registerParticipantsToEvent(List<String> emails, TripModel trip) {
    List<ParticipantModel> res = emails
      .stream()
      .map(email -> new ParticipantModel(email, trip))
      .toList();

    this.repository.saveAll(res);
  };
  @CacheEvict(key="#tripId" , value="Participant")
  public ParticipantCreatedResponse registerParticipantToEvent(String email, UUID tripId) {

    TripModel tripModel = this.tripRepository.findById(tripId).orElse(null);

    ParticipantModel newParticipantModel = new ParticipantModel(email, tripModel);
    this.repository.save(newParticipantModel);

    return new ParticipantCreatedResponse(newParticipantModel.getId());
  }


  public void triggerConfirmationEmailToParticipants(UUID tripId) {};
  public void triggerConfirmationEmailToParticipant(String email) {};
  @Cacheable("Participant")
  public List<ParticipantData> getAllParticipantsByTripId(UUID tripId) {
    return this
      .repository
      .findByTripId(tripId)
      .stream()
      .map(participant -> new ParticipantData(participant.getId(), participant.getName(), participant.getEmail(), participant.getIsConfirmed()))
      .toList();
  }

  @Caching(evict = {
    @CacheEvict(value = "Participant", allEntries = true),
    @CacheEvict(value = "ParticipantItem", key = "#id")
  })
  public boolean confirmParticipant(UUID id, ParticipantRequestPayload payload) {
    var optional = this.repository.findById(id);
    if (optional.isEmpty()) return false;
    var entity = optional.get();
    entity.setIsConfirmed(true);
    entity.setName(payload.name());
    this.repository.save(entity);
    return true;
  }

  
  @Caching(evict = {
          @CacheEvict(value = "Participant", allEntries = true),
          @CacheEvict(value = "ParticipantItem", key = "#participantId")
  })
  public boolean updateParticipantConfirmationStatus(UUID participantId, Boolean isConfirmed) {
    var participantOpt = this.repository.findById(participantId);
    if (participantOpt.isEmpty()) return false;
    var participant = participantOpt.get();
    participant.setIsConfirmed(isConfirmed);
    this.repository.save(participant);
    return true;
  }

  @Caching(evict = {
          @CacheEvict(value = "Participant", allEntries = true),
          @CacheEvict(value = "ParticipantItem", key = "#participantId")
  })
  public boolean updateParticipantDetails(UUID participantId, ParticipantRequestPayload payload) {
    var participantOpt = this.repository.findById(participantId);
    if (participantOpt.isEmpty()) return false;
    var participant = participantOpt.get();
    participant.setName(payload.name());
    participant.setEmail(payload.email());
    this.repository.save(participant);
    return true;
  }

  @Caching(evict = {
          @CacheEvict(value = "Participant", allEntries = true),
          @CacheEvict(value = "ParticipantItem", key = "#participantId")
  })
  public boolean deleteParticipantsById(UUID participantId) {
    if (!this.repository.existsById(participantId)) return false;
    this.repository.deleteById(participantId);
    return true;
  }

  @Cacheable(value = "ParticipantItem", key = "#id")
  public java.util.Optional<ParticipantData> getParticipantById(UUID id) {
    return this.repository.findById(id)
            .map(p -> new ParticipantData(p.getId(), p.getName(), p.getEmail(), p.getIsConfirmed()));
  }


}

