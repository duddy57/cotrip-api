package com.cotrip.links;

import com.cotrip.links.DTO.LinkData;
import com.cotrip.links.DTO.LinkRequestPayload;
import com.cotrip.links.DTO.LinkResponse;
import com.cotrip.trip.TripModel;
import com.cotrip.trip.TripRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.Caching;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
public class LinkService {

  @Autowired
  private LinkRepository repository;
  @Autowired
  private TripRepository tripRepository;

  @CacheEvict(key="#tripId" , value="Link")
  public void createLink(LinkRequestPayload payload, UUID tripId) {

      var trip = tripRepository.findById(tripId).orElse(null);

      LinkModel newLink = new LinkModel(payload.title(), payload.url(), trip);

    this
      .repository
      .save(newLink);
  }

  @Cacheable("Link")
  public List<LinkData> getAllLinkForTrip(UUID trip_id) {
    return this
      .repository
      .findByTripId(trip_id)
      .stream()
      .map(activity -> new LinkData(activity.getId(), activity.getTitle(), activity.getUrl()))
      .toList();
  }

  @Cacheable(value = "LinkItem", key = "#tripId.toString() + ':' + #linkId.toString()")
  public Optional<LinkData> getLinkByIdForTrip(UUID tripId, UUID linkId) {
    Optional<LinkModel> optional = this.repository.findById(linkId);
    if (optional.isEmpty()) return Optional.empty();
    LinkModel entity = optional.get();
    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return Optional.empty();
    }
    return Optional.of(new LinkData(entity.getId(), entity.getTitle(), entity.getUrl()));
  }

  @Caching(evict = {
    @CacheEvict(key = "#tripId", value = "Link"),
    @CacheEvict(key = "#tripId.toString() + ':' + #linkId.toString()", value = "LinkItem")
  })
  public Optional<LinkData> updateLink(UUID tripId, UUID linkId, LinkRequestPayload payload) {
    Optional<LinkModel> optional = this.repository.findById(linkId);
    if (optional.isEmpty()) return Optional.empty();

    LinkModel entity = optional.get();
    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return Optional.empty();
    }

    if (payload.title() != null) entity.setTitle(payload.title());
    if (payload.url() != null) entity.setUrl(payload.url());

    this.repository.save(entity);

    return Optional.of(new LinkData(entity.getId(), entity.getTitle(), entity.getUrl()));
  }

  @Caching(evict = {
    @CacheEvict(key = "#tripId", value = "Link"),
    @CacheEvict(key = "#tripId.toString() + ':' + #linkId.toString()", value = "LinkItem")
  })
  public boolean deleteLink(UUID tripId, UUID linkId) {
    Optional<LinkModel> optional = this.repository.findById(linkId);
    if (optional.isEmpty()) return false;

    LinkModel entity = optional.get();
    if (entity.getTrip() == null || entity.getTrip().getId() == null || !entity.getTrip().getId().equals(tripId)) {
      return false;
    }

    this.repository.delete(entity);
    return true;
  }
}
