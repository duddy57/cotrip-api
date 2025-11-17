package com.cotrip.trip;

import com.cotrip.participant.ParticipantService;
import com.cotrip.trip.DTO.TripRequestPayload;
import com.cotrip.trip.DTO.TripGetDTO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
public class TripService {

    @Autowired
    TripRepository tripRepository;

    @Autowired
    ParticipantService participantService;

    public TripGetDTO CreateTrip(TripRequestPayload tripRequestPayload) {

        TripModel newTripModel = new TripModel(tripRequestPayload);

        this.tripRepository.save(newTripModel);

        this.participantService.registerParticipantsToEvent(tripRequestPayload.emails_to_invite(), newTripModel);

        return new TripGetDTO(
                newTripModel.getId()
                ,newTripModel.getDestination()
                ,newTripModel.getStartAt()
                ,newTripModel.getEndAt()
                , newTripModel.getIsConfirmed()
        );

    }

    public List<TripGetDTO> getTrips() {

        List<TripModel> trips = this.tripRepository.findAll();
        List <TripGetDTO> tripGetDTOs = new ArrayList<>();
        for(TripModel newTripModel : trips) {

            tripGetDTOs.add( new TripGetDTO(
                     newTripModel.getId()
                    ,newTripModel.getDestination()
                    ,newTripModel.getStartAt()
                    ,newTripModel.getEndAt()
                    ,newTripModel.getIsConfirmed()
            ));
        }

        return tripGetDTOs;

    }

    public Optional<TripGetDTO>  getTripById(UUID tripId) {

        Optional<TripModel> newTripModel = this.tripRepository.findById(tripId);

        return newTripModel.map(tripModel -> new TripGetDTO(
                tripModel.getId()
                , tripModel.getDestination()
                , tripModel.getStartAt()
                , tripModel.getEndAt()
                , tripModel.getIsConfirmed()
        ));

    }

    public List<Optional<TripGetDTO>> getTripsByOwnerEmail(String ownerEmail) {

        List<TripModel> tripModels = this.tripRepository.findByOwnerEmail(ownerEmail);
        List<Optional<TripGetDTO>> tripGetDTOs = new ArrayList<>();

        for(TripModel tripModel : tripModels) {
            tripGetDTOs.add( Optional.of(new TripGetDTO(
                    tripModel.getId()
                    , tripModel.getDestination()
                    , tripModel.getStartAt()
                    , tripModel.getEndAt()
                    , tripModel.getIsConfirmed()
            )));
        }

        return tripGetDTOs;

    }

    public Optional<TripGetDTO> UpdateTrip(UUID tripId, TripRequestPayload payload) {

        Optional<TripModel>  optionalTripModel = tripRepository.findById(tripId);

        if (optionalTripModel.isPresent()) {

            var rawTripModel = optionalTripModel.get();

            rawTripModel.setDestination(payload.destination());
            rawTripModel.setStartAt(LocalDateTime.parse(payload.start_at(), DateTimeFormatter.ISO_DATE_TIME));
            rawTripModel.setEndAt(LocalDateTime.parse(payload.end_at(), DateTimeFormatter.ISO_DATE_TIME));

            this.tripRepository.save(rawTripModel);

        }
        return optionalTripModel.map(tripModel -> new TripGetDTO(
                tripModel.getId()
                , tripModel.getDestination()
                , tripModel.getStartAt()
                , tripModel.getEndAt()
                , tripModel.getIsConfirmed()
        ));
    }

    public Optional<TripGetDTO> ConfirmTrip(UUID tripId) {

        Optional<TripModel> optionalTripModel = this.tripRepository.findById(tripId);

        if(optionalTripModel.isPresent()) {
            TripModel rawTripModel = optionalTripModel.get();
            rawTripModel.setIsConfirmed(true);
            this.tripRepository.save(rawTripModel);
        }

        return optionalTripModel.map(tripModel -> new TripGetDTO(
                tripModel.getId()
                , tripModel.getDestination()
                , tripModel.getStartAt()
                , tripModel.getEndAt()
                , tripModel.getIsConfirmed()
        ));

    }


}
