/**
 * Copyright (C) 2015, 2016  Green Screens Ltd.
 */
class SmartCardUtil {

	static toArray(string_data) {
		return new TextEncoder("utf-8").encode(string_data);
	}

	static toBinary(base64_string) {
		return Uint8Array.from(atob(base64_string), c => c.charCodeAt(0));
	}

	static toHex(binary_data) {
		return Array.prototype.map.call(binary_data, x => ('00' + x.toString(16)).slice(-2)).join('');
	}

	static fromHex(hex_string) {
		return new Uint8Array(hex_string.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));
	}

	static toBase64(uint8_array) {
		return btoa(String.fromCharCode.apply(null, uint8_array));
	}
}

class SmartCard {

	static log(data) {

		let el = document.querySelector('textarea');
		if (typeof data === 'string') {
			el.value = el.value + '\n' + data;
			return;
		}

		let Tlv = data.data ? data.data.Tlv ||null : null;

		if (Array.isArray(Tlv)) {
			Tlv.every(tlv => {
				if (tlv.Tag === "53") {
					let hash = SmartCardUtil.toBase64(SmartCardUtil.fromHex(tlv.Data));
					window.open(location.origin + `/asn1#${hash}`, '_blank');
				}
				return true;
			});
		}

		el.value = el.value + '\n' + JSON.stringify(data);
	}

	static get Type() {
		return parseInt(document.querySelector('#type').value);
	}

	static get Cls() {
		return parseInt(document.querySelector('#cls').value);
	}

	static get Ins() {
		return parseInt(document.querySelector('#ins').value);
	}

	static get P1() {
		return parseInt(document.querySelector('#p1').value);
	}

	static get P2() {
		return parseInt(document.querySelector('#p2').value);
	}

	static get Data() {

		let val = document.querySelector('#data').value.trim();

		if (val) {
			val = SmartCardUtil.fromHex(val);
			val = SmartCardUtil.toBase64(val);
		}

		return val;
	}

	static get Oid() {
		return parseInt(document.querySelector('#oid').value);
	}

	static get Pin() {

		let val = document.querySelector('#pin').value;
		val = SmartCardUtil.toArray(val.padEnd(8, '\0'));

		val.forEach(function(item, i) {
			if (item == 0) val[i] = 255;
		});

		return SmartCardUtil.toBase64(val);
	}

	static get command() {
		return {
			type: SmartCard.Type,
			cla: SmartCard.Cls,
			ins: SmartCard.Ins,
			p1: SmartCard.P1,
			p2: SmartCard.P2,
			data: SmartCard.Data,
		};
	}

	static async _recieve(url, opt) {

		let response = await fetch(location.origin + url, opt || {});

		let data = await response.text();

		if (response.status == 200) {
			try {
				data = JSON.parse(data)
			} catch (e) {
				console.log(e);
			}
		}

		SmartCard.log(data);
		return data;
	}

	static async list() {
		SmartCard.log('List Smart Cards...');
		return SmartCard._recieve('/list');
	}

	static async connect(id = 0) {
		SmartCard.log('Connecting...');
		return SmartCard._recieve('/connect?id=' + id);
	}

	static async disconnect() {
		SmartCard.log('Disconnecting...');
		return SmartCard._recieve('/disconnect');
	}

	static async send() {

		SmartCard.log('Sending...');

		let obj = SmartCard.command;

		let args = {
			method: 'PUT',
			credentials: 'same-origin',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(obj)
		};

		return SmartCard._recieve('/request', args);
	}

	static async authenticate() {
		SmartCard.log('Authenticating...');
		let pin = SmartCard.Pin
		return SmartCard._recieve(`/pin?id=${pin}`);
	}

	static async valid() {
		SmartCard.log('Valid...');
		return SmartCard._recieve('/valid');
	}

	static async version() {
		SmartCard.log('Version...');
		return SmartCard._recieve('/version');
	}

	static async biometric() {
		SmartCard.log('Biometric...');
		return SmartCard._recieve('/bio');
	}

	static async dob() {
		SmartCard.log('Discovery Object...');
		return SmartCard._recieve('/dob');
	}

	static async oidIdentifier() {
		SmartCard.log('Object identifier...');
		let oid = SmartCard.Oid
		return SmartCard._recieve(`/oid?id=${oid}`);
	}

	static clear() {
		document.querySelector('textarea').value = '';
	}

}
