import {EventEmitter} from "./event-emitter.js";

const ukDefaultLatLng = { lat: 55.378051, lng: -3.435973 }
const mapBoxAccessToken = 'pk.eyJ1Ijoic2lsa3Rvd24tc29mdHdhcmUiLCJhIjoiY2w0ZXY4dWk4MDIxYTNmbWZiMmh2dG1ueiJ9.sbEpaf3_hJqph0NjxQ3XAg';
const defaultMapZoom = 4;
const localMapZoom = 14;

export class MapboxContainer {
    /**
     * Ctor
     * @param selector
     * @param eventEmitter
     */
    constructor(selector, eventEmitter) {
        if (!selector) {
            throw new Error('selector string is empty');
        }

        if (!eventEmitter || !(eventEmitter instanceof EventEmitter)) {
            throw new Error('event emitter is invalid');
        }

        this.container = document.querySelector(selector);
        this.eventEmitter = eventEmitter;
    }

    /**
     * initialises
     */
    init() {
        mapboxgl.accessToken = mapBoxAccessToken;

        this.map = new mapboxgl.Map({
            container: 'map', // container ID
            style: 'mapbox://styles/mapbox/streets-v12', // style URL
            center: [ukDefaultLatLng.lng, ukDefaultLatLng.lat], // starting position [lng, lat]
            zoom: defaultMapZoom, // starting zoom
        });
    }

    /**
     * Binds our map events
     */
    bind() {
        this.eventEmitter.on('ON_RETRIEVE_POSTCODE_LOCATION', (data) => this.#updateMapLocation(data));
        this.eventEmitter.on("ON_USER_RESET_MAP", () => this.#resetMap())
    }

    /**
     * Resets the map back to the default view when first loaded
     */
    #resetMap() {
        if (!this.marker) {
            return;
        }

        this.marker.remove();

        this.map.setZoom(defaultMapZoom);
        this.map.setCenter([ukDefaultLatLng.lng, ukDefaultLatLng.lat]);
    }

    /**
     * Sets the map location based on the lng/lat provided
     * @param data the postcode location
     */
    #updateMapLocation(data) {
        const { lat, lng } = data;

        /* set the map centre and zoom in*/
        this.map.setCenter({lat, lng});
        this.map.setZoom(localMapZoom);

        /* add a marker */
        this.marker = new mapboxgl.Marker()
            .setLngLat([lng, lat])
            .addTo(this.map);
    }
}