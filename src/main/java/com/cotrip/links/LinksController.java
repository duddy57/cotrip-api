package com.cotrip.links;

import com.cotrip.links.DTO.LinkData;
import com.cotrip.links.DTO.LinkRequestPayload;
import com.cotrip.trip.DTO.TripGetDTO;
import com.cotrip.trip.TripService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@RestController
@RequestMapping("/links")
public class LinksController {

    private final TripService tripService;
    private final LinkService linkService;

    @Autowired
    public LinksController(TripService tripService, LinkService linkService) {
        this.tripService = tripService;
        this.linkService = linkService;
    }

    @PostMapping("/{tripId}/new-link")
    public ResponseEntity<String> createNewLink(@PathVariable UUID tripId, @RequestBody LinkRequestPayload payload) {
        Optional<TripGetDTO> trip = tripService.getTripById(tripId);

        if(trip.isEmpty()) {
            return ResponseEntity.notFound().build();
        }

        this.linkService.createLink(payload, trip.get().id());


        return ResponseEntity.ok().build();
    }

    @GetMapping("/{tripId}/links")
    public ResponseEntity<List<LinkData>> getAllLinks(@PathVariable UUID tripId) {
        List<LinkData> res = this.linkService.getAllLinkForTrip(tripId);
        return ResponseEntity.ok(res);

    }

    @GetMapping("/{tripId}/links/{linkId}")
    public ResponseEntity<LinkData> getLinkById(@PathVariable UUID tripId, @PathVariable UUID linkId) {
        var res = this.linkService.getLinkByIdForTrip(tripId, linkId);
        return res.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @PutMapping("/{tripId}/links/{linkId}")
    public ResponseEntity<LinkData> updateLink(@PathVariable UUID tripId,
                                               @PathVariable UUID linkId,
                                               @RequestBody LinkRequestPayload payload) {
        var updated = this.linkService.updateLink(tripId, linkId, payload);
        return updated.map(ResponseEntity::ok).orElseGet(() -> ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{tripId}/links/{linkId}")
    public ResponseEntity<Void> deleteLink(@PathVariable UUID tripId, @PathVariable UUID linkId) {
        boolean deleted = this.linkService.deleteLink(tripId, linkId);
        if (!deleted) return ResponseEntity.notFound().build();
        return ResponseEntity.noContent().build();
    }
}
