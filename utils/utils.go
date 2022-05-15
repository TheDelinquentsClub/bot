package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/kingultron99/tdcbot/logger"
)

// BToMb converts byte values to MegaBytes
func BToMb(b uint64) string {
	return fmt.Sprintf("%.1f", float64(b/1024)/1024)
}

// GetDurationString returns a time duration for use with uptime or other related uses
func GetDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

func ConvertTickToDuration(ticks int) string {
	var years = ticks / 630720000
	ticks = ticks % 630720000
	var days = ticks / 1728000
	ticks = ticks % 1728000
	var hours = ticks / 72000
	ticks = ticks % 72000
	var minutes = ticks / 1200
	ticks = ticks % 1200
	var seconds = ticks / 20
	ticks = ticks % 20
	return fmt.Sprintf("%v:%v:%v:%v:%v", years, days, hours, minutes, seconds)
}

// DefaultColour is the default discord.color to use in embeds
// DiscordGreen is the colour to be used in signifying a success message, or something good
// DiscordRed is the colour to be used in signifying an error message, or something bad
var (
	DefaultColour discord.Color = 0xA3BCF9
	DiscordGreen  discord.Color = 0x379A57
	DiscordBlue   discord.Color = 0x5865F2
	DiscordRed    discord.Color = 0xDF3E41
)

func MustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(env)
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}

type Player struct {
	Username string `json:"name"`
	UUID     string `json:"id"`
}

// GetUUID returns the UUID tied with to the username provided
func GetUUID(username string) string {
	username = strings.ReplaceAll(username, "\"", "")
	url := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%v", username)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var user = new(Player)
	err = json.Unmarshal(body, &user)
	return user.UUID

}

// GetUsername returns the username of a player from the provided UUID
func GetUsername(UUID string) string {
	UUID = strings.ReplaceAll(UUID, "\"", "")
	url := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%v", UUID)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var user = new(Player)
	err = json.Unmarshal(body, &user)
	return user.Username
}

type PlayerNames struct {
	Name    string `json:"name"`
	Changed int64  `json:"changedToAt,omitempty"`
}

// GetNamesFromUsername returns all the usernames the specified player has had using the provided username
func GetNamesFromUsername(username string) string {
	uuid := GetUUID(username)
	return GetNamesFromUUID(uuid)
}

// GetNamesFromUUID returns all the usernames the specified player has had using the provided UUID
func GetNamesFromUUID(uuid string) string {
	uuid = strings.ReplaceAll(uuid, "\"", "")
	url := fmt.Sprintf("https://api.mojang.com/user/profiles/%v/names", uuid)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}

	var names = new([]PlayerNames)
	err = json.Unmarshal(body, &names)

	var (
		nameArray []string
		res       string
	)

	for _, playerNames := range *names {
		//convert ms to s
		changed := fmt.Sprintf("<t:%v:R>\n", playerNames.Changed/1000)
		if playerNames.Changed == 0 {
			changed = "Accounts first username!\n"
		}
		nameArray = append(nameArray, fmt.Sprintf("%v — %v", playerNames.Name, changed))
	}

	if len(nameArray) >= 6 {
		concat := strings.Builder{}
		concat.WriteString(strings.Join(nameArray[:2], ""))
		concat.WriteString("...\n")
		concat.WriteString(strings.Join(nameArray[len(nameArray)-3:len(nameArray)-1], ""))
		res = concat.String()
	} else {
		res = strings.Join(nameArray, "")
	}
	return res
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandString(length int) string {
	return StringWithCharset(length, charset)
}

type LocaleTypes struct {
	Advancements_adventure_adventuring_time_title                             string `json:"advancements.Adventure.Adventuring_Time.Title"`
	Advancements_adventure_adventuring_time_description                       string `json:"advancements.adventure.adventuring_time.description"`
	Advancements_adventure_arbalistic_title                                   string `json:"advancements.adventure.arbalistic.title"`
	Advancements_adventure_arbalistic_description                             string `json:"advancements.adventure.arbalistic.description"`
	Advancements_adventure_bullseye_title                                     string `json:"advancements.adventure.bullseye.title"`
	Advancements_adventure_bullseye_description                               string `json:"advancements.adventure.bullseye.description"`
	Advancements_adventure_fall_from_world_height_title                       string `json:"advancements.adventure.fall_from_world_height.title"`
	Advancements_adventure_fall_from_world_height_description                 string `json:"advancements.adventure.fall_from_world_height.description"`
	Advancements_adventure_walk_on_powder_snow_with_leather_boots_title       string `json:"advancements.adventure.walk_on_powder_snow_with_leather_boots.title"`
	Advancements_adventure_walk_on_powder_snow_with_leather_boots_description string `json:"advancements.adventure.walk_on_powder_snow_with_leather_boots.description"`
	Advancements_adventure_lightning_rod_with_villager_no_fire_title          string `json:"advancements.adventure.lightning_rod_with_villager_no_fire.title"`
	Advancements_adventure_lightning_rod_with_villager_no_fire_description    string `json:"advancements.adventure.lightning_rod_with_villager_no_fire.description"`
	Advancements_adventure_spyglass_at_parrot_title                           string `json:"advancements.adventure.spyglass_at_parrot.title"`
	Advancements_adventure_spyglass_at_parrot_description                     string `json:"advancements.adventure.spyglass_at_parrot.description"`
	Advancements_adventure_spyglass_at_ghast_title                            string `json:"advancements.adventure.spyglass_at_ghast.title"`
	Advancements_adventure_spyglass_at_ghast_description                      string `json:"advancements.adventure.spyglass_at_ghast.description"`
	Advancements_adventure_spyglass_at_dragon_title                           string `json:"advancements.adventure.spyglass_at_dragon.title"`
	Advancements_adventure_spyglass_at_dragon_description                     string `json:"advancements.adventure.spyglass_at_dragon.description"`
	Advancements_adventure_hero_of_the_village_title                          string `json:"advancements.adventure.hero_of_the._illage.title"`
	Advancements_adventure_hero_of_the_village_description                    string `json:"advancements.adventure.hero_of_the_village.description"`
	Advancements_adventure_honey_block_slide_title                            string `json:"advancements.adventure.honey_block_slide.title"`
	Advancements_adventure_honey_block_slide_description                      string `json:"advancements.adventure.honey_block_slide.description"`
	Advancements_adventure_kill_all_mobs_title                                string `json:"advancements.adventure.kill_all_mobs.title"`
	Advancements_adventure_kill_all_mobs_description                          string `json:"advancements.adventure.kill_all_mobs.description"`
	Advancements_adventure_kill_a_mob_title                                   string `json:"advancements.adventure.kill_a_mob.title"`
	Advancements_adventure_kill_a_mob_description                             string `json:"advancements.adventure.kill_a_mob.description"`
	Advancements_adventure_ol_betsy_title                                     string `json:"advancements.adventure.ol_betsy.title"`
	Advancements_adventure_ol_betsy_description                               string `json:"advancements.adventure.ol_betsy.description"`
	Advancements_adventure_play_jukebox_in_meadows_title                      string `json:"advancements.adventure.play_jukebox_in_meadows.title"`
	Advancements_adventure_play_jukebox_in_meadows_description                string `json:"advancements.adventure.play_jukebox_in_meadows.description"`
	Advancements_adventure_shoot_arrow_title                                  string `json:"advancements.adventure.shoot_arrow.title"`
	Advancements_adventure_shoot_arrow_description                            string `json:"advancements.adventure.shoot_arrow.description"`
	Advancements_adventure_sleep_in_bed_title                                 string `json:"advancements.adventure.sleep_in_bed.title"`
	Advancements_adventure_sleep_in_bed_description                           string `json:"advancements.adventure.sleep_in_bed.description"`
	Advancements_adventure_sniper_duel_title                                  string `json:"advancements.adventure.sniper_duel.title"`
	Advancements_adventure_sniper_duel_description                            string `json:"advancements.adventure.sniper_duel.description"`
	Advancements_adventure_summon_iron_golem_title                            string `json:"advancements.adventure.summon_iron_golem.title"`
	Advancements_adventure_summon_iron_golem_description                      string `json:"advancements.adventure.summon_iron_golem.description"`
	Advancements_adventure_totem_of_undying_title                             string `json:"advancements.adventure.totem_of_undying.title"`
	Advancements_adventure_totem_of_undying_description                       string `json:"advancements.adventure.totem_of_undying.description"`
	Advancements_adventure_trade_title                                        string `json:"advancements.adventure.trade.title"`
	Advancements_adventure_trade_description                                  string `json:"advancements.adventure.trade.description"`
	Advancements_adventure_trade_at_world_height_title                        string `json:"advancements.adventure.trade_at_world_height.title"`
	Advancements_adventure_trade_at_world_height_description                  string `json:"advancements.adventure.trade_at_world_height.description"`
	Advancements_adventure_throw_trident_title                                string `json:"advancements.adventure.throw_trident.title"`
	Advancements_adventure_throw_trident_description                          string `json:"advancements.adventure.throw_trident.description"`
	Advancements_adventure_two_birds_one_arrow_title                          string `json:"advancements.adventure.two_birds_one_arrow.title"`
	Advancements_adventure_two_birds_one_arrow_description                    string `json:"advancements.adventure.two_birds_one_arrow.description"`
	Advancements_adventure_very_very_frightening_title                        string `json:"advancements.adventure.very_very_frightening.title"`
	Advancements_adventure_very_very_frightening_description                  string `json:"advancements.adventure.very_very_frightening.description"`
	Advancements_adventure_voluntary_exile_title                              string `json:"advancements.adventure.voluntary_exile.title"`
	Advancements_adventure_voluntary_exile_description                        string `json:"advancements.adventure.voluntary_exile.description"`
	Advancements_adventure_whos_the_pillager_now_title                        string `json:"advancements.adventure.whos_the_pillager_now.title"`
	Advancements_adventure_whos_the_pillager_now_description                  string `json:"advancements.adventure.whos_he_pillager_now.description"`
	Advancements_husbandry_breed_an_animal_title                              string `json:"advancements.husbandry.breed_an_animal.title"`
	Advancements_husbandry_breed_an_animal_description                        string `json:"advancements.husbandry.breed_an_animal.description"`
	Advancements_husbandry_fishy_business_title                               string `json:"advancements.husbandry.fishy_business.title"`
	Advancements_husbandry_fishy_business_description                         string `json:"advancements.husbandry.fishy_business.description"`
	Advancements_husbandry_make_a_sign_glow_title                             string `json:"advancements.husbandry.make_a_sign_glow.title"`
	Advancements_husbandry_make_a_sign_glow_description                       string `json:"advancements.husbandry.make_a_sign_glow.description"`
	Advancements_husbandry_ride_a_boat_with_a_goat_title                      string `json:"advancements.husbandry.ride_a_boat_with_a_goat.title"`
	Advancements_husbandry_ride_a_boat_with_a_goat_description                string `json:"advancements.husbandry.ride_a_boat_with_a_goat.description"`
	Advancements_husbandry_tactical_fishing_title                             string `json:"advancements.husbandry.tactical_fishing.title"`
	Advancements_husbandry_tactical_fishing_description                       string `json:"advancements.husbandry.tactical_fishing.description"`
	Advancements_husbandry_axolotl_in_a_bucket_title                          string `json:"advancements.husbandry.axolotl_in_a_bucket.title"`
	Advancements_husbandry_axolotl_in_a_bucket_description                    string `json:"advancements.husbandry.axolotl_in_a_bucket.description"`
	Advancements_husbandry_kill_axolotl_target_title                          string `json:"advancements.husbandry.kill_axolotl_target.title"`
	Advancements_husbandry_kill_axolotl_target_description                    string `json:"advancements.husbandry.kill_axolotl_target.description"`
	Advancements_husbandry_breed_all_animals_title                            string `json:"advancements.husbandry.breed_all_animals.title"`
	Advancements_husbandry_breed_all_animals_description                      string `json:"advancements.husbandry.breed_all_animals.description"`
	Advancements_husbandry_tame_an_animal_title                               string `json:"advancements.husbandry.tame_an_animal.title"`
	Advancements_husbandry_tame_an_animal_description                         string `json:"advancements.husbandry.tame_an_animal.description"`
	Advancements_husbandry_plant_seed_title                                   string `json:"advancements.husbandry.plant_seed.title"`
	Advancements_husbandry_plant_seed_description                             string `json:"advancements.husbandry.plant_seed.description"`
	Advancements_husbandry_netherite_hoe_title                                string `json:"advancements.husbandry.netherite_hoe.title"`
	Advancements_husbandry_netherite_hoe_description                          string `json:"advancements.husbandry.netherite_hoe.description"`
	Advancements_husbandry_balanced_diet_title                                string `json:"advancements.husbandry.balanced_diet.title"`
	Advancements_husbandry_balanced_diet_description                          string `json:"advancements.husbandry.balanced_diet.description"`
	Advancements_husbandry_complete_catalogue_title                           string `json:"advancements.husbandry.complete_catalogue.title"`
	Advancements_husbandry_complete_catalogue_description                     string `json:"advancements.husbandry.complete_catalogue.description"`
	Advancements_husbandry_safely_harvest_honey_title                         string `json:"advancements.husbandry.safely_harvest_honey.title"`
	Advancements_husbandry_safely_harvest_honey_description                   string `json:"advancements.husbandry.safely_harvest_honey.description"`
	Advancements_husbandry_silk_touch_nest_title                              string `json:"advancements.husbandry.silk_touch_nest.title"`
	Advancements_husbandry_silk_touch_nest_description                        string `json:"advancements.husbandry.silk_touch_nest.description"`
	Advancements_husbandry_wax_on_title                                       string `json:"advancements.husbandry.wax_on.title"`
	Advancements_husbandry_wax_on_description                                 string `json:"advancements.husbandry.wax_on.description"`
	Advancements_husbandry_wax_off_title                                      string `json:"advancements.husbandry.wax_off.title"`
	Advancements_husbandry_wax_off_description                                string `json:"advancements.husbandry.wax_off.description"`
	Advancements_end_dragon_breath_title                                      string `json:"advancements.end.dragon_breath.title"`
	Advancements_end_dragon_breath_description                                string `json:"advancements.end.dragon_breath.description"`
	Advancements_end_dragon_egg_title                                         string `json:"advancements.end.dragon_egg.title"`
	Advancements_end_dragon_egg_description                                   string `json:"advancements.end.dragon_egg.description"`
	Advancements_end_elytra_title                                             string `json:"advancements.end.elytra.title"`
	Advancements_end_elytra_description                                       string `json:"advancements.end.elytra.description"`
	Advancements_end_enter_end_gateway_title                                  string `json:"advancements.end.enter_end_gateway.title"`
	Advancements_end_enter_end_gateway_description                            string `json:"advancements.end.enter_end_gateway.description"`
	Advancements_end_find_end_city_title                                      string `json:"advancements.end.find_end_city.title"`
	Advancements_end_find_end_city_description                                string `json:"advancements.end.find_end_city.description"`
	Advancements_end_kill_dragon_title                                        string `json:"advancements.end.kill_dragon.title"`
	Advancements_end_kill_dragon_description                                  string `json:"advancements.end.kill_dragon.description"`
	Advancements_end_levitate_title                                           string `json:"advancements.end.levitate.title"`
	Advancements_end_levitate_description                                     string `json:"advancements.end.levitate.description"`
	Advancements_end_respawn_dragon_title                                     string `json:"advancements.end.respawn_dragon.title"`
	Advancements_end_respawn_dragon_description                               string `json:"advancements.end.respawn_dragon.description"`
	Advancements_nether_brew_potion_description                               string `json:"advancements.nether.brew_potion.description"`
	Advancements_nether_all_potions_title                                     string `json:"advancements.nether.all_potions.title"`
	Advancements_nether_all_potions_description                               string `json:"advancements.nether.all_potions.description"`
	Advancements_nether_all_effects_title                                     string `json:"advancements.nether.all_effects.title"`
	Advancements_nether_all_effects_description                               string `json:"advancements.nether.all_effects.description"`
	Advancements_nether_create_beacon_title                                   string `json:"advancements.nether.create_beacon.title"`
	Advancements_nether_create_beacon_description                             string `json:"advancements.nether.create_beacon.description"`
	Advancements_nether_create_full_beacon_title                              string `json:"advancements.nether.create_full_beacon.title"`
	Advancements_nether_create_full_beacon_description                        string `json:"advancements.nether.create_full_beacon.description"`
	Advancements_nether_find_fortress_title                                   string `json:"advancements.nether.find_fortress.title"`
	Advancements_nether_find_fortress_description                             string `json:"advancements.nether.find_fortress.description"`
	Advancements_nether_get_wither_skull_title                                string `json:"advancements.nether.get_wither_skull.title"`
	Advancements_nether_get_wither_skull_description                          string `json:"advancements.nether.get_wither_skull.description"`
	Advancements_nether_obtain_blaze_rod_title                                string `json:"advancements.nether.obtain_blaze_rod.title"`
	Advancements_nether_obtain_blaze_rod_description                          string `json:"advancements.nether.obtain_blaze_rod.description"`
	Advancements_nether_return_to_sender_title                                string `json:"advancements.nether.return_to_sender.title"`
	Advancements_nether_return_to_sender_description                          string `json:"advancements.nether.return_to_sender.description"`
	Advancements_nether_summon_wither_description                             string `json:"advancements.nether.summon_wither.description"`
	Advancements_nether_fast_travel_title                                     string `json:"advancements.nether.fast_travel.title"`
	Advancements_nether_fast_travel_description                               string `json:"advancements.nether.fast_travel.description"`
	Advancements_nether_uneasy_alliance_title                                 string `json:"advancements.nether.uneasy_alliance.title"`
	Advancements_nether_uneasy_alliance_description                           string `json:"advancements.nether.uneasy_alliance.description"`
	Advancements_nether_obtain_ancient_debris_title                           string `json:"advancements.nether.obtain_ancient_debris.title"`
	Advancements_nether_obtain_ancient_debris_description                     string `json:"advancements.nether.obtain_ancient_debris.description"`
	Advancements_nether_netherite_armor_title                                 string `json:"advancements.nether.netherite_armor.title"`
	Advancements_nether_netherite_armor_description                           string `json:"advancements.nether.netherite_armor.description"`
	Advancements_nether_use_lodestone_title                                   string `json:"advancements.nether.use_lodestone.title"`
	Advancements_nether_use_lodestone_description                             string `json:"advancements.nether.use_lodestone.description"`
	Advancements_nether_obtain_crying_obsidian_title                          string `json:"advancements.nether.obtain_crying_obsidian.title"`
	Advancements_nether_obtain_crying_obsidian_description                    string `json:"advancements.nether.obtain_crying_obsidian.description"`
	Advancements_nether_charge_respawn_anchor_title                           string `json:"advancements.nether.charge_respawn_anchor.title"`
	Advancements_nether_charge_respawn_anchor_description                     string `json:"advancements.nether.charge_respawn_anchor.description"`
	Advancements_nether_ride_strider_title                                    string `json:"advancements.nether.ride_strider.title"`
	Advancements_nether_ride_strider_description                              string `json:"advancements.nether.ride_strider.description"`
	Advancements_nether_ride_strider_in_overworld_lava_title                  string `json:"advancements.nether.ride_strider_in_overworld_lava.title"`
	Advancements_nether_ride_strider_in_overworld_lava_description            string `json:"advancements.nether.ride_strider_in_overworld_lava.description"`
	Advancements_nether_explore_nether_title                                  string `json:"advancements.nether.explore_nether.title"`
	Advancements_nether_explore_nether_description                            string `json:"advancements.nether.explore_nether.description"`
	Advancements_nether_find_bastion_title                                    string `json:"advancements.nether.find_bastion.title"`
	Advancements_nether_find_bastion_description                              string `json:"advancements.nether.find_bastion.description"`
	Advancements_nether_loot_bastion_title                                    string `json:"advancements.nether.loot_bastion.title"`
	Advancements_nether_loot_bastion_description                              string `json:"advancements.nether.loot_bastion.description"`
	Advancements_nether_distract_piglin_title                                 string `json:"advancements.nether.distract_piglin.title"`
	Advancements_nether_distract_piglin_description                           string `json:"advancements.nether.distract_piglin.description"`
	Advancements_story_cure_zombie_villager_title                             string `json:"advancements.story.cure_zombie_villager.title"`
	Advancements_story_cure_zombie_villager_description                       string `json:"advancements.story.cure_zombie_villager.description"`
	Advancements_story_deflect_arrow_title                                    string `json:"advancements.story.deflect_arrow.title"`
	Advancements_story_deflect_arrow_description                              string `json:"advancements.story.deflect_arrow.description"`
	Advancements_story_enchant_item_title                                     string `json:"advancements.story.enchant_item.title"`
	Advancements_story_enchant_item_description                               string `json:"advancements.story.enchant_item.description"`
	Advancements_story_enter_the_end_title                                    string `json:"advancements.story.enter_the_end.title"`
	Advancements_story_enter_the_end_description                              string `json:"advancements.story.enter_the_end.description"`
	Advancements_story_enter_the_nether_title                                 string `json:"advancements.story.enter_the_nether.title"`
	Advancements_story_enter_the_nether_description                           string `json:"advancements.story.enter_the_nether.description"`
	Advancements_story_follow_ender_eye_title                                 string `json:"advancements.story.follow_ender_eye.title"`
	Advancements_story_follow_ender_eye_description                           string `json:"advancements.story.follow_ender_eye.description"`
	Advancements_story_form_obsidian_title                                    string `json:"advancements.story.form_obsidian.title"`
	Advancements_story_form_obsidian_description                              string `json:"advancements.story.form_obsidian.description"`
	Advancements_story_iron_tools_title                                       string `json:"advancements.story.iron_tools.title"`
	Advancements_story_iron_tools_description                                 string `json:"advancements.story.iron_tools.description"`
	Advancements_story_lava_bucket_title                                      string `json:"advancements.story.lava_bucket.title"`
	Advancements_story_lava_bucket_description                                string `json:"advancements.story.lava_bucket.description"`
	Advancements_story_mine_diamond_title                                     string `json:"advancements.story.mine_diamond.title"`
	Advancements_story_mine_diamond_description                               string `json:"advancements.story.mine_diamond.description"`
	Advancements_story_mine_stone_title                                       string `json:"advancements.story.mine_stone.title"`
	Advancements_story_mine_stone_description                                 string `json:"advancements.story.mine_stone.description"`
	Advancements_story_obtain_armor_title                                     string `json:"advancements.story.obtain_armor.title"`
	Advancements_story_obtain_armor_description                               string `json:"advancements.story.obtain_armor.description"`
	Advancements_story_shiny_gear_title                                       string `json:"advancements.story.shiny_gear.title"`
	Advancements_story_shiny_gear_description                                 string `json:"advancements.story.shiny_gear.description"`
	Advancements_story_smelt_iron_title                                       string `json:"advancements.story.smelt_iron.title"`
	Advancements_story_smelt_iron_description                                 string `json:"advancements.story.smelt_iron.description"`
	Advancements_story_upgrade_tools_title                                    string `json:"advancements.story.upgrade_tools.title"`
	Advancements_story_upgrade_tools_description                              string `json:"advancements.story.upgrade_tools.description"`
}

var Locale = new(LocaleTypes)

func LoadLocale() {
	localeFile, _ := os.Open("./assets/mc_en_US.json")
	_ = json.NewDecoder(localeFile).Decode(&Locale)
}

func GetLocale(locale string, localeType string) string {
	r := reflect.ValueOf(Locale)
	res := reflect.Indirect(r).FieldByName(strings.Title(locale) + "_" + localeType).String()
	return res
}
