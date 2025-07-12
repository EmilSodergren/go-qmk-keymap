// qmk-keymap-format.json
#include QMK_KEYBOARD_H
#include <stdio.h>

enum lily_layers {
    _COLEMAK,
    _SIGNS,
    _NUMBERS
};

const uint16_t PROGMEM keymaps[][MATRIX_ROWS][MATRIX_COLS] = {
    //# .-----------------------------------------------.                                             .-----------------------------------------------.
    //# |   `   | B_VAI | B_VAD | B_TOG |   ⍉   |   ⍉   |                                             |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# |  TAB  |   Q   |   W   |   F   |   P   |   G   |                                             |   J   |   L   |   U   |   Y   |   ;   |   ⍉   |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# | MCR ⎋ |   A   |   R   | S(⎇)  | T(⌘)  |   D   |-------.                             .-------|   H   | N(⌘)  | E(⎇)  |   I   |   O   | QUOTE |
    //# |-------+-------+-------+-------+-------+-------|   ⍉   |                             |   ⍉   |-------+-------+-------+-------+-------+-------|
    //# | SHIFT |   Z   |   X   |   C   |   V   |   B   |-------|                             |-------|   K   |   M   |   ,   |   .   |   /   |   ⍉   |
    //# .-------+-------+-------+-------+-------+-------/       /                             \       \-------+-------+-------+-------+-------+-------.
    //#                      |   ⌘   | MO(4) |   ⎇   | /   ⎇   /                               \   ⌫   \  | SL(1) | RCTRL |   ⍉   |
    //#                      |       |       |       |/       /                                 \       \ |       |       |       |
    //#                      .-------+-------+-------+-------.                                    .------++-------+-------+-------.
  [_COLEMAK] = LAYOUT(
    KC_GRAVE,               RGB_VAI,  RGB_VAD,  RGB_TOG,     KC_NO,       KC_NO,                        KC_NO,  KC_NO,       KC_NO,       KC_NO,    KC_NO,     KC_NO,
    KC_TAB,                 KC_Q,     KC_W,     KC_F,        KC_P,        KC_G,                         KC_J,   KC_L,        KC_U,        KC_Y,     KC_SCOLON, KC_NO,    // apa is my friend
    RESET_LAYER_AND_ESCAPE, KC_A,     KC_R,     OPT_T(KC_S), CMD_T(KC_T), KC_D,                         KC_H,   CMD_T(KC_N), OPT_T(KC_E), KC_I,     KC_O,      KC_QUOTE,
    OSM(KC_LSHIFT),         KC_Z,     KC_X,     KC_C,        KC_V,        KC_B,    KC_NO,    KC_NO,     KC_K,   KC_M,        KC_COMMA,    KC_DOT,   KC_SLASH,  KC_NO,
                                                KC_LGUI,     MO(4),       KC_LALT, KC_LALT,  KC_BSPACE, OSL(1), KC_RCTRL,    KC_NO
  ),

    //# .-----------------------------------------------.                                             .-----------------------------------------------.
    //# | ARENT | KC_1) |   ⍉   |   ⍉   |   ⍉   |   ⍉   |                                             |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# |       |   !   |   @   |   #   |   $   |   %   |                                             |   ^   |   [   |   ]   |   ⍉   |   ⍉   |   ⍉   |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# |       |   <   |   _   |   =   |   >   |   ~   |-------.                             .-------|   ⍉   |   (   |   )   |   ⍉   |   ⏎   |       |
    //# |-------+-------+-------+-------+-------+-------|   ⍉   |                             |   ⍉   |-------+-------+-------+-------+-------+-------|
    //# |       |   `   |   ⍉   |   ,   |   .   |   ⍉   |-------|                             |-------|   \   |   {   |   }   |   *   | TO(3) |       |
    //# .-------+-------+-------+-------+-------+-------/       /                             \       \-------+-------+-------+-------+-------+-------.
    //#                      |       |   -   |   +   | /       /                               \   ⌦   \  | TO(2) |       |   ⍉   |
    //#                      |       |       |       |/       /                                 \       \ |       |       |       |
    //#                      .-------+-------+-------+-------.                                    .------++-------+-------+-------.
  [_SIGNS] = LAYOUT(
    C_TRANSPARENT,  LGUI(KC_1), KC_NO,    KC_NO,          KC_NO,    KC_NO,                              KC_NO,     KC_NO,          KC_NO,       KC_NO,    KC_NO,    KC_NO,
    KC_TRANSPARENT, KC_EXLM,    KC_AT,    KC_HASH,        KC_DLR,   KC_PERC,                            KC_CIRC,   KC_LBRACKET,    KC_RBRACKET, KC_NO,    KC_NO,    KC_NO,
    KC_TRANSPARENT, KC_LABK,    KC_UNDS,  KC_EQUAL,       KC_RABK,  KC_TILD,                            KC_NO,     KC_LPRN,        KC_RPRN,     KC_NO,    KC_ENTER, KC_TRANSPARENT,
    KC_TRANSPARENT, KC_GRAVE,   KC_NO,    KC_COMMA,       KC_DOT,   KC_NO,   KC_NO,          KC_NO,     KC_BSLASH, KC_LCBR,        KC_RCBR,     KC_ASTR,  TO(3),    KC_TRANSPARENT,
                                          KC_TRANSPARENT, KC_MINUS, KC_PLUS, KC_TRANSPARENT, KC_DELETE, TO(2),     KC_TRANSPARENT, KC_NO
  ),

    //# .-----------------------------------------------.                                             .-----------------------------------------------.
    //# |       |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |                                             |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |   ⍉   |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# |       |       |   /   |   -   |   +   |       |                                             |   ⍉   |   7   |   8   |   9   |       |       |
    //# |-------+-------+-------+-------+-------+-------|                                             |-------+-------+-------+-------+-------+-------|
    //# |       |       |   *   |   1   |   0   |   =   |-------.                             .-------|   ⍉   |   4   |   5   |   6   |   ⏎   |       |
    //# |-------+-------+-------+-------+-------+-------|   ⍉   |                             |   ⍉   |-------+-------+-------+-------+-------+-------|
    //# |       |       |   _   |   ,   |   .   |       |-------|                             |-------|   0   |   1   |   2   |   3   |   /   |       |
    //# .-------+-------+-------+-------+-------+-------/       /                             \       \-------+-------+-------+-------+-------+-------.
    //#                      |       | TO(0) |   ⌘   | /       /                               \       \  |       |       |       |
    //#                      |       |       |       |/       /                                 \       \ |       |       |       |
    //#                      .-------+-------+-------+-------.                                    .------++-------+-------+-------.
  [_NUMBERS] = LAYOUT(
    KC_TRANSPARENT, KC_NO,          KC_NO,          KC_NO,          KC_NO,   KC_NO,                                          KC_NO,          KC_NO,          KC_NO,          KC_NO,    KC_NO,          KC_NO,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_SLASH,       KC_MINUS,       KC_PLUS, KC_TRANSPARENT,                                 KC_NO,          KC_7,           KC_8,           KC_9,     KC_TRANSPARENT, KC_TRANSPARENT,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_KP_ASTERISK, KC_1,           KC_0,    KC_EQUAL,                                       KC_NO,          KC_4,           KC_5,           KC_6,     KC_ENTER,       KC_TRANSPARENT,
    KC_TRANSPARENT, KC_TRANSPARENT, KC_UNDS,        KC_COMMA,       KC_DOT,  KC_TRANSPARENT, KC_NO,          KC_NO,          KC_0,           KC_1,           KC_2,           KC_3,     KC_SLASH,       KC_TRANSPARENT,
                                                    KC_TRANSPARENT, TO(0),   KC_LGUI,        KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT, KC_TRANSPARENT
    )
};
