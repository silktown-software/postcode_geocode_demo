import {EventEmitter} from './event-emitter.js';
import {PostcodeForm} from './postcode-form.js';
import {MapboxContainer} from './mapbox-container.js';

document.addEventListener('DOMContentLoaded', () => {
    const eventEmitter = new EventEmitter();

    const postcodeForm = new PostcodeForm('postalcode-form', eventEmitter);
    postcodeForm.init();
    postcodeForm.bind();

    const mapboxContainer = new MapboxContainer('#mapbox-container', eventEmitter);
    mapboxContainer.init();
    mapboxContainer.bind();
});